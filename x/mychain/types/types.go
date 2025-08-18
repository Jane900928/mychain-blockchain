package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// User 用户结构体
type User struct {
	Address   string    `json:"address" yaml:"address"`
	Name      string    `json:"name" yaml:"name"`
	Email     string    `json:"email" yaml:"email"`
	Balance   sdk.Coins `json:"balance" yaml:"balance"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

// NewUser 创建新用户
func NewUser(address, name, email string) User {
	return User{
		Address:   address,
		Name:      name,
		Email:     email,
		Balance:   sdk.NewCoins(),
		CreatedAt: time.Now(),
	}
}

// Transaction 交易结构体
type Transaction struct {
	Hash        string    `json:"hash" yaml:"hash"`
	From        string    `json:"from" yaml:"from"`
	To          string    `json:"to" yaml:"to"`
	Amount      sdk.Coins `json:"amount" yaml:"amount"`
	Fee         sdk.Coins `json:"fee" yaml:"fee"`
	BlockHeight int64     `json:"block_height" yaml:"block_height"`
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	Status      string    `json:"status" yaml:"status"`
}

// NewTransaction 创建新交易
func NewTransaction(hash, from, to string, amount, fee sdk.Coins, blockHeight int64) Transaction {
	return Transaction{
		Hash:        hash,
		From:        from,
		To:          to,
		Amount:      amount,
		Fee:         fee,
		BlockHeight: blockHeight,
		Timestamp:   time.Now(),
		Status:      "confirmed",
	}
}

// Miner 矿工结构体
type Miner struct {
	Address      string    `json:"address" yaml:"address"`
	Description  string    `json:"description" yaml:"description"`
	Commission   sdk.Dec   `json:"commission" yaml:"commission"`
	TotalRewards sdk.Coins `json:"total_rewards" yaml:"total_rewards"`
	BlocksMined  int64     `json:"blocks_mined" yaml:"blocks_mined"`
	Status       string    `json:"status" yaml:"status"`
	RegisteredAt time.Time `json:"registered_at" yaml:"registered_at"`
	LastActiveAt time.Time `json:"last_active_at" yaml:"last_active_at"`
}

// NewMiner 创建新矿工
func NewMiner(address, description string, commission sdk.Dec) Miner {
	return Miner{
		Address:      address,
		Description:  description,
		Commission:   commission,
		TotalRewards: sdk.NewCoins(),
		BlocksMined:  0,
		Status:       "active",
		RegisteredAt: time.Now(),
		LastActiveAt: time.Now(),
	}
}

// BlockInfo 区块信息结构体
type BlockInfo struct {
	Height    int64     `json:"height" yaml:"height"`
	Hash      string    `json:"hash" yaml:"hash"`
	PrevHash  string    `json:"prev_hash" yaml:"prev_hash"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Miner     string    `json:"miner" yaml:"miner"`
	TxCount   int       `json:"tx_count" yaml:"tx_count"`
	TotalFees sdk.Coins `json:"total_fees" yaml:"total_fees"`
	Reward    sdk.Coins `json:"reward" yaml:"reward"`
	Size      int64     `json:"size" yaml:"size"`
	GasUsed   int64     `json:"gas_used" yaml:"gas_used"`
	GasLimit  int64     `json:"gas_limit" yaml:"gas_limit"`
}

// NewBlockInfo 创建新区块信息
func NewBlockInfo(height int64, hash, prevHash, miner string, txCount int) BlockInfo {
	return BlockInfo{
		Height:    height,
		Hash:      hash,
		PrevHash:  prevHash,
		Timestamp: time.Now(),
		Miner:     miner,
		TxCount:   txCount,
		TotalFees: sdk.NewCoins(),
		Reward:    sdk.NewCoins(sdk.NewCoin("mychain", sdk.NewInt(100))), // 默认挖矿奖励
		Size:      0,
		GasUsed:   0,
		GasLimit:  10000000,
	}
}

// TokenSupply 代币供应量结构体
type TokenSupply struct {
	TotalSupply       sdk.Coins `json:"total_supply" yaml:"total_supply"`
	CirculatingSupply sdk.Coins `json:"circulating_supply" yaml:"circulating_supply"`
	MaxSupply         sdk.Coins `json:"max_supply" yaml:"max_supply"`
	MintedToday       sdk.Coins `json:"minted_today" yaml:"minted_today"`
	LastMintTime      time.Time `json:"last_mint_time" yaml:"last_mint_time"`
}

// NewTokenSupply 创建新的代币供应量
func NewTokenSupply() TokenSupply {
	maxSupply := sdk.NewCoins(sdk.NewCoin("mychain", sdk.NewInt(1000000000))) // 10亿代币上限
	return TokenSupply{
		TotalSupply:       sdk.NewCoins(),
		CirculatingSupply: sdk.NewCoins(),
		MaxSupply:         maxSupply,
		MintedToday:       sdk.NewCoins(),
		LastMintTime:      time.Now(),
	}
}

// Params 模块参数
type Params struct {
	MaxSupply      sdk.Coins `json:"max_supply" yaml:"max_supply"`
	MiningReward   sdk.Coins `json:"mining_reward" yaml:"mining_reward"`
	TransactionFee sdk.Coins `json:"transaction_fee" yaml:"transaction_fee"`
	MinCommission  sdk.Dec   `json:"min_commission" yaml:"min_commission"`
	MaxCommission  sdk.Dec   `json:"max_commission" yaml:"max_commission"`
	BlockTime      int64     `json:"block_time" yaml:"block_time"` // 秒
	MaxValidators  int       `json:"max_validators" yaml:"max_validators"`
}

// DefaultParams 返回默认参数
func DefaultParams() Params {
	return Params{
		MaxSupply:      sdk.NewCoins(sdk.NewCoin("mychain", sdk.NewInt(1000000000))),
		MiningReward:   sdk.NewCoins(sdk.NewCoin("mychain", sdk.NewInt(100))),
		TransactionFee: sdk.NewCoins(sdk.NewCoin("mychain", sdk.NewInt(1))),
		MinCommission:  sdk.NewDecWithPrec(0, 2),  // 0%
		MaxCommission:  sdk.NewDecWithPrec(20, 2), // 20%
		BlockTime:      6,                         // 6秒出块
		MaxValidators:  100,                       // 最多100个验证者
	}
}

// Validate 验证参数
func (p Params) Validate() error {
	if !p.MaxSupply.IsValid() {
		return fmt.Errorf("invalid max supply: %s", p.MaxSupply)
	}
	
	if !p.MiningReward.IsValid() {
		return fmt.Errorf("invalid mining reward: %s", p.MiningReward)
	}
	
	if !p.TransactionFee.IsValid() {
		return fmt.Errorf("invalid transaction fee: %s", p.TransactionFee)
	}
	
	if p.MinCommission.IsNegative() || p.MinCommission.GT(sdk.OneDec()) {
		return fmt.Errorf("invalid min commission: %s", p.MinCommission)
	}
	
	if p.MaxCommission.IsNegative() || p.MaxCommission.GT(sdk.OneDec()) {
		return fmt.Errorf("invalid max commission: %s", p.MaxCommission)
	}
	
	if p.MinCommission.GT(p.MaxCommission) {
		return fmt.Errorf("min commission cannot be greater than max commission")
	}
	
	if p.BlockTime <= 0 {
		return fmt.Errorf("block time must be positive")
	}
	
	if p.MaxValidators <= 0 {
		return fmt.Errorf("max validators must be positive")
	}
	
	return nil
}

// String 返回参数的字符串表示
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  MaxSupply:         %s
  MiningReward:      %s
  TransactionFee:    %s
  MinCommission:     %s
  MaxCommission:     %s
  BlockTime:         %d
  MaxValidators:     %d`, 
		p.MaxSupply, p.MiningReward, p.TransactionFee, p.MinCommission, p.MaxCommission, p.BlockTime, p.MaxValidators)
}

// ParamKeyTable 返回参数键表
func ParamKeyTable() interface{} {
	// 简化实现，实际应该返回 paramtypes.KeyTable
	return nil
}

// GenesisState 创世状态
type GenesisState struct {
	Params       Params        `json:"params" yaml:"params"`
	Users        []User        `json:"users" yaml:"users"`
	Miners       []Miner       `json:"miners" yaml:"miners"`
	Transactions []Transaction `json:"transactions" yaml:"transactions"`
	TokenSupply  TokenSupply   `json:"token_supply" yaml:"token_supply"`
}

// DefaultGenesis 返回默认创世状态
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		Users:        []User{},
		Miners:       []Miner{},
		Transactions: []Transaction{},
		TokenSupply:  NewTokenSupply(),
	}
}

// Validate 验证创世状态
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
