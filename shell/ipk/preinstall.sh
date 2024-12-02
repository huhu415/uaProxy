# 检查是否为Linux系统
if [ "$(uname)" != "Linux" ]; then
    echo "This script must be run on Linux"
    exit 1
fi


# 检查系统启动方式
if [ ! -f /etc/rc.common ] || [ "$(cat /proc/1/comm)" != "procd" ]; then
    echo "This package is designed for OpenWrt systems only"
    exit 1
fi

# 检查必要的系统命令
for cmd in iptables grep; do
    if ! command -v $cmd >/dev/null 2>&1; then
        echo "Required command '$cmd' not found"
        exit 1
    fi
done
