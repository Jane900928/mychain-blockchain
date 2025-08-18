package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	// 自定义模块
	"github.com/Jane900928/mychain-blockchain/x/mychain"
	mychainkeeper "github.com/Jane900928/mychain-blockchain/x/mychain/keeper"
	mychaintypes "github.com/Jane900928/mychain-blockchain/x/mychain/types"
)

const (
	AppName      = "MyChain"
	Bech32Prefix = "mychain"
)

var (
	// DefaultNodeHome 默认节点主目录
	DefaultNodeHome string

	// Bech32PrefixAccAddr 定义主网的帐户地址 bech32 前缀
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub 定义主网的账户公钥 bech32 前缀
	Bech32PrefixAccPub = Bech32Prefix + "pub"
	// Bech32PrefixValAddr 定义主网的验证器运营商地址 bech32 前缀
	Bech32PrefixValAddr = Bech32Prefix + "valoper"
	// Bech32PrefixValPub 定义主网的验证器运营商公钥 bech32 前缀
	Bech32PrefixValPub = Bech32Prefix + "valoperpub"
	// Bech32PrefixConsAddr 定义主网的验证器节点地址 bech32 前缀
	Bech32PrefixConsAddr = Bech32Prefix + "valcons"
	// Bech32PrefixConsPub 定义主网的验证器节点公钥 bech32 前缀
	Bech32PrefixConsPub = Bech32Prefix + "valconspub"

	// ModuleBasics 定义应用程序使用的模块
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		mychain.AppModuleBasic{}, // 自定义模块
	)

	// 模块账户权限
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		mychaintypes.ModuleName:        {authtypes.Minter},
	}
)

// MyChainApp 扩展 ABCI 应用程序
type MyChainApp struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// Keys 用于访问子存储
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// Keepers
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper stakingkeeper.Keeper
	ParamsKeeper  paramskeeper.Keeper

	// 自定义模块的 Keeper
	MyChainKeeper mychainkeeper.Keeper

	// 模块管理器
	mm *module.Manager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".mychain")
}

// NewMyChainApp 返回对初始化的 MyChain 的引用
func NewMyChainApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig EncodingConfig,
	appOpts interface{},
	baseAppOptions ...func(*baseapp.BaseApp),
) *MyChainApp {

	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		paramstypes.StoreKey, mychaintypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys()

	app := &MyChainApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// 初始化 params keeper 和子空间
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// 设置特定的参数子空间
	authSubspace := app.ParamsKeeper.Subspace(authtypes.ModuleName)
	bankSubspace := app.ParamsKeeper.Subspace(banktypes.ModuleName)
	stakingSubspace := app.ParamsKeeper.Subspace(stakingtypes.ModuleName)
	mychainSubspace := app.ParamsKeeper.Subspace(mychaintypes.ModuleName)

	// 添加 keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authSubspace, authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, bankSubspace, app.ModuleAccountAddrs(),
	)
	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, stakingSubspace,
	)

	// 初始化自定义模块 keeper
	app.MyChainKeeper = mychainkeeper.NewKeeper(
		appCodec,
		keys[mychaintypes.StoreKey],
		mychainSubspace,
		app.BankKeeper,
	)

	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		params.NewAppModule(app.ParamsKeeper),
		mychain.NewAppModule(appCodec, app.MyChainKeeper), // 自定义模块
	)

	app.mm.SetOrderInitGenesis(
		authtypes.ModuleName, banktypes.StoreKey, stakingtypes.ModuleName,
		paramstypes.ModuleName, mychaintypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	// 初始化存储
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// 初始化 BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  nil,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// 应用程序编码配置
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig 创建 EncodingConfig
func MakeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

func (app *MyChainApp) BeginBlocker(ctx sdk.Context, req interface{}) interface{} {
	return app.mm.BeginBlock(ctx, req)
}

func (app *MyChainApp) EndBlocker(ctx sdk.Context, req interface{}) interface{} {
	return app.mm.EndBlock(ctx, req)
}

func (app *MyChainApp) InitChainer(ctx sdk.Context, req interface{}) interface{} {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal([]byte("{}"), &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *MyChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *MyChainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// initParamsKeeper 初始化参数 keeper 和它的子空间
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(mychaintypes.ModuleName)
	return paramsKeeper
}
