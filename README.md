# MyChain 区块链项目

🚀 一个基于 Cosmos SDK 开发的完整区块链解决方案

## 功能特性

- ✅ **代币管理** - 铸造、转账、供应量控制
- ✅ **用户系统** - 创建用户、管理账户
- ✅ **矿工机制** - 注册矿工、奖励分发
- ✅ **区块链浏览器** - Web 界面查看链上数据
- ✅ **CosmJS 集成** - 完整的前端客户端库

## 快速开始

### 1. 环境要求

- Go 1.21+
- Node.js 16+
- Git

### 2. 编译项目

```bash
# 克隆项目
git clone https://github.com/Jane900928/mychain-blockchain.git
cd mychain-blockchain

# 安装依赖
go mod tidy

# 编译
go build -o build/mychaind ./cmd/mychaind
```

### 3. 初始化和启动

```bash
# 给脚本执行权限
chmod +x scripts/deploy.sh

# 初始化区块链
./scripts/deploy.sh init

# 启动节点
./scripts/deploy.sh start
```

### 4. 使用区块链浏览器

打开 `web/explorer/index.html` 文件，或使用本地服务器：

```bash
cd web/explorer
python3 -m http.server 8080
# 访问 http://localhost:8080
```

## 项目结构

```
mychain-blockchain/
├── cmd/mychaind/           # 主程序入口
├── x/mychain/              # 自定义模块
│   ├── types/              # 类型定义
│   ├── keeper/             # 状态管理
│   └── client/cli/         # CLI 命令
├── app/                    # 应用程序逻辑
├── web/explorer/           # 区块链浏览器
├── client/                 # CosmJS 客户端
└── scripts/                # 部署脚本
```

## 主要命令

### 用户管理
```bash
# 创建用户
mychaind tx mychain create-user "Alice" "alice@example.com" --from validator

# 查询用户
mychaind query mychain user mychain1xxxxx
```

### 代币操作
```bash
# 铸造代币
mychaind tx mychain mint 1000000mychain --from validator

# 转账
mychaind tx mychain transfer mychain1receiver 1000mychain --from sender
```

### 矿工功能
```bash
# 注册矿工
mychaind tx mychain register-miner "My Pool" "0.05" --from miner

# 查询矿工
mychaind query mychain miners
```

## CosmJS 客户端

```typescript
import { MyChainClient } from './client/cosmjs-client';

const client = new MyChainClient();
await client.connect();

const address = await client.connectWithMnemonic(mnemonic);
const txHash = await client.createUser("Alice", "alice@example.com", address);
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
