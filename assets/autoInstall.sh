#!/bin/bash

ARCH_TYPE=""
get_arch_type() {
    local ARCH=$(uname -m)

    case $ARCH in
        i386|i686)
            ARCH_TYPE="386"
            ;;
        x86_64)
            ARCH_TYPE="amd64"
            ;;
        armv7*|armv6*|armv8*)
            ARCH_TYPE="arm"
            ;;
        aarch64)
            ARCH_TYPE="arm64"
            ;;
        loongarch64)
            ARCH_TYPE="loong64"
            ;;
        mips)
            # ARCH_TYPE="mips" 大端
            ARCH_TYPE="mipsle"
            ;;
        mips64)
            # ARCH_TYPE="mips64" 大端
            ARCH_TYPE="mips64le"
            ;;
        ppc64)
            # ARCH_TYPE="ppc64" 大端
            ARCH_TYPE="ppc64le"
            ;;
        riscv64)
            ARCH_TYPE="riscv64"
            ;;
        *)
            ARCH_TYPE=$ARCH
            return 1
            ;;
    esac

    return 0
}


# 检查是否为root用户
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root"
   exit 1
fi

if [ "$(uname)" != "Linux" ]; then
    echo "This script must be run on Linux"
    exit 1
fi

# 检查是否为支持的架构
if get_arch_type; then
    echo "Detected architecture type: $ARCH_TYPE"
else
    echo "$ARCH_TYPE arch is not supported"
    exit 1
fi

# 1. 开启 IP 转发
echo "正在开启 IP 转发..."
if ! grep -q "net.ipv4.ip_forward=1" /etc/sysctl.conf; then
    echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
fi
sysctl -p

# 2. 下载并安装 uaProxy
echo "正在下载 uaProxy.tar.gz..."
# 这里需要根据实际情况修改下载链接
UaProxy_Name="uaProxy_linux_${ARCH_TYPE}"
DOWNLOAD_URL="https://github.com/huhu415/uaProxy/releases/latest/download/${UaProxy_Name}.tar.gz"

# 下载
if ! wget -q "$DOWNLOAD_URL"; then
    echo "下载失败,请检查网络连接或下载链接是否正确"
    exit 1
fi

tar -xzf $UaProxy_Name.tar.gz
cd $UaProxy_Name

mv uaProxy /usr/sbin/
chmod +x /usr/sbin/uaProxy

mv assets/uaProxy-openwrt /etc/init.d/
chmod +x /etc/init.d/uaProxy-openwrt

# 3. 配置 iptables 规则
echo "正在配置 iptables 规则..."
iptables -t nat -N uaProxy
iptables -t nat -A uaProxy -d 192.168.0.0/16 -j RETURN
iptables -t nat -A uaProxy -p tcp -j RETURN -m mark --mark 0xff
iptables -t nat -A uaProxy -p tcp -j REDIRECT --to-ports 12345
iptables -t nat -A PREROUTING -p tcp -j uaProxy
iptables -t nat -A OUTPUT -p tcp -j uaProxy

# 启用并启动服务
echo "正在启用服务..."
/etc/init.d/uaProxy-openwrt enable
/etc/init.d/uaProxy-openwrt start

cd ..
rm -rf $UaProxy_Name

echo "安装完成！"
echo "可以使用 '/etc/init.d/uaProxy-openwrt {start|stop|restart|status}' 来控制服务"
echo "可以使用 'logread | grep uaProxy' 查看运行日志"
