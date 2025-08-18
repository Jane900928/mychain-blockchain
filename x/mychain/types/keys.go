package types

const (
	// ModuleName 定义模块名称
	ModuleName = "mychain"

	// StoreKey 定义模块的主存储键
	StoreKey = ModuleName

	// RouterKey 定义模块的消息路由键
	RouterKey = ModuleName

	// QuerierRoute 定义模块的查询路由
	QuerierRoute = ModuleName
)

// 存储键
var (
	// TokenSupplyKey 代币供应量键
	TokenSupplyKey = []byte{0x01}
	
	// UserKeyPrefix 用户键前缀
	UserKeyPrefix = []byte{0x02}
	
	// TransactionKeyPrefix 交易键前缀
	TransactionKeyPrefix = []byte{0x03}
	
	// MinerKeyPrefix 矿工键前缀
	MinerKeyPrefix = []byte{0x04}
	
	// BlockInfoKeyPrefix 区块信息键前缀
	BlockInfoKeyPrefix = []byte{0x05}
)

// GetUserKey 返回用户的存储键
func GetUserKey(address string) []byte {
	return append(UserKeyPrefix, []byte(address)...)
}

// GetTransactionKey 返回交易的存储键
func GetTransactionKey(txHash string) []byte {
	return append(TransactionKeyPrefix, []byte(txHash)...)
}

// GetMinerKey 返回矿工的存储键
func GetMinerKey(address string) []byte {
	return append(MinerKeyPrefix, []byte(address)...)
}

// GetBlockInfoKey 返回区块信息的存储键
func GetBlockInfoKey(height int64) []byte {
	heightBytes := make([]byte, 8)
	// 简化的高度转换
	for i := 0; i < 8; i++ {
		heightBytes[i] = byte(height >> (8 * (7 - i)))
	}
	return append(BlockInfoKeyPrefix, heightBytes...)
}
