#!/bin/bash

# MyChain 区块链部署和启动脚本
# 使用方法: ./deploy.sh [init|start|stop|restart|clean]

set -e

# 配置变量
CHAIN_ID="mychain-1"
MONIKER="mychain-node"
KEYRING_BACKEND="test"
HOME_DIR="$HOME/.mychain"
BINARY_NAME="mychaind"
VALIDATOR_KEY="validator"
USER_KEY="user"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印彩色消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查二进制文件是否存在
check_binary() {
    if ! command -v $BINARY_NAME &> /dev/null; then
        print_error "$BINARY_NAME 二进制文件未找到，请先编译项目"
        exit 1
    fi
    print_message "找到 $BINARY_NAME 二进制文件"
}

# 初始化链
init_chain() {
    print_step "初始化 MyChain 区块链..."
    
    # 清理旧数据
    if [ -d "$HOME_DIR" ]; then
        print_warning "删除现有的链数据目录: $HOME_DIR"
        rm -rf "$HOME_DIR"
    fi
    
    # 初始化节点
    print_message "初始化节点配置..."
    $BINARY_NAME init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR
    
    # 创建验证者密钥
    print_message "创建验证者密钥..."
    $BINARY_NAME keys add $VALIDATOR_KEY --keyring-backend $KEYRING_BACKEND --home $HOME_DIR
    
    # 创建用户密钥
    print_message "创建用户密钥..."
    $BINARY_NAME keys add $USER_KEY --keyring-backend $KEYRING_BACKEND --home $HOME_DIR
    
    # 获取地址
    VALIDATOR_ADDR=$($BINARY_NAME keys show $VALIDATOR_KEY -a --keyring-backend $KEYRING_BACKEND --home $HOME_DIR)
    USER_ADDR=$($BINARY_NAME keys show $USER_KEY -a --keyring-backend $KEYRING_BACKEND --home $HOME_DIR)
    
    print_message "验证者地址: $VALIDATOR_ADDR"
    print_message "用户地址: $USER_ADDR"
    
    # 将验证者地址添加到创世状态
    print_message "配置创世状态..."
    $BINARY_NAME add-genesis-account $VALIDATOR_ADDR 100000000000000000000000000mychain --home $HOME_DIR
    $BINARY_NAME add-genesis-account $USER_ADDR 1000000000000000000000000mychain --home $HOME_DIR
    
    # 创建创世交易
    print_message "创建创世交易..."
    $BINARY_NAME gentx $VALIDATOR_KEY 1000000000000000000000000mychain \
        --chain-id $CHAIN_ID \
        --keyring-backend $KEYRING_BACKEND \
        --home $HOME_DIR
    
    # 收集创世交易
    print_message "收集创世交易..."
    $BINARY_NAME collect-gentxs --home $HOME_DIR
    
    # 验证创世文件
    print_message "验证创世文件..."
    $BINARY_NAME validate-genesis --home $HOME_DIR
    
    # 配置节点
    configure_node
    
    print_message "MyChain 初始化完成！"
    print_message "验证者地址: $VALIDATOR_ADDR"
    print_message "用户地址: $USER_ADDR"
}

# 配置节点
configure_node() {
    print_step "配置节点参数..."
    
    # 配置文件路径
    CONFIG_FILE="$HOME_DIR/config/config.toml"
    APP_FILE="$HOME_DIR/config/app.toml"
    
    # 修改配置文件
    if [ -f "$CONFIG_FILE" ]; then
        print_message "配置 config.toml..."
        
        # 启用 RPC 服务器
        sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $CONFIG_FILE
        
        # 设置出块时间
        sed -i 's/timeout_commit = "5s"/timeout_commit = "6s"/g' $CONFIG_FILE
    fi
    
    if [ -f "$APP_FILE" ]; then
        print_message "配置 app.toml..."
        
        # 启用 REST API
        sed -i 's/enable = false/enable = true/g' $APP_FILE
        sed -i 's/address = "tcp:\/\/0.0.0.0:1317"/address = "tcp:\/\/0.0.0.0:1317"/g' $APP_FILE
        
        # 配置最小 gas 价格
        sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001mychain"/g' $APP_FILE
    fi
    
    print_message "节点配置完成"
}

# 启动链
start_chain() {
    print_step "启动 MyChain 区块链..."
    
    if [ ! -d "$HOME_DIR" ]; then
        print_error "链未初始化，请先运行: $0 init"
        exit 1
    fi
    
    print_message "启动节点..."
    print_message "RPC 端点: http://localhost:26657"
    print_message "REST API: http://localhost:1317"
    
    # 启动节点
    $BINARY_NAME start --home $HOME_DIR
}

# 停止链
stop_chain() {
    print_step "停止 MyChain 区块链..."
    
    # 查找并杀死进程
    PIDS=$(pgrep -f $BINARY_NAME || true)
    if [ -n "$PIDS" ]; then
        print_message "停止进程: $PIDS"
        kill $PIDS
        sleep 2
        
        # 强制杀死进程（如果还在运行）
        PIDS=$(pgrep -f $BINARY_NAME || true)
        if [ -n "$PIDS" ]; then
            print_warning "强制停止进程: $PIDS"
            kill -9 $PIDS
        fi
    else
        print_message "没有找到运行中的 $BINARY_NAME 进程"
    fi
}

# 重启链
restart_chain() {
    print_step "重启 MyChain 区块链..."
    stop_chain
    sleep 2
    start_chain
}

# 清理数据
clean_chain() {
    print_step "清理 MyChain 数据..."
    
    stop_chain
    
    if [ -d "$HOME_DIR" ]; then
        print_warning "删除链数据目录: $HOME_DIR"
        rm -rf "$HOME_DIR"
        print_message "清理完成"
    else
        print_message "没有找到数据目录"
    fi
}

# 显示状态
show_status() {
    print_step "MyChain 状态信息..."
    
    # 检查进程
    PIDS=$(pgrep -f $BINARY_NAME || true)
    if [ -n "$PIDS" ]; then
        print_message "状态: 运行中 (PID: $PIDS)"
    else
        print_message "状态: 未运行"
    fi
    
    # 检查数据目录
    if [ -d "$HOME_DIR" ]; then
        print_message "数据目录: $HOME_DIR (已初始化)"
    else
        print_message "数据目录: 未初始化"
    fi
    
    # 检查端口
    print_message "检查端口状态..."
    if command -v netstat &> /dev/null; then
        netstat -tuln | grep -E ":(26657|1317)" || print_message "相关端口未监听"
    fi
}

# 显示帮助信息
show_help() {
    echo "MyChain 区块链管理脚本"
    echo ""
    echo "用法: $0 [COMMAND]"
    echo ""
    echo "命令:"
    echo "  init     - 初始化区块链"
    echo "  start    - 启动区块链"
    echo "  stop     - 停止区块链"
    echo "  restart  - 重启区块链"
    echo "  clean    - 清理所有数据"
    echo "  status   - 显示状态信息"
    echo "  help     - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 init      # 初始化新的区块链"
    echo "  $0 start     # 启动区块链服务"
    echo "  $0 stop      # 停止区块链服务"
    echo "  $0 clean     # 删除所有数据并重新开始"
}

# 主函数
main() {
    case "${1:-}" in
        init)
            check_binary
            init_chain
            ;;
        start)
            check_binary
            start_chain
            ;;
        stop)
            stop_chain
            ;;
        restart)
            check_binary
            restart_chain
            ;;
        clean)
            clean_chain
            ;;
        status)
            show_status
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: ${1:-}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
