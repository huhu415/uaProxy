#!/bin/sh

ARCH_TYPE=""
get_arch_type() {
    local ARCH=$(uname -m)

    # 检测大小端
    local mipstype=$(echo -n I | hexdump -o 2>/dev/null | awk '{ print substr($2,6,1); exit}')

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
            # 如果mipstype=0是大端，如果是1则是小端
            if [ "$mipstype" = "0" ]; then
                ARCH_TYPE="mips"    # 大端
            else
                ARCH_TYPE="mipsle"  # 小端
            fi
            ;;
        mips64)
            # 如果mipstype=0是大端，如果是1则是小端
            if [ "$mipstype" = "0" ]; then
                ARCH_TYPE="mips64"    # 大端
            else
                ARCH_TYPE="mips64le"  # 小端
            fi
            ;;
        ppc64)
            # 如果mipstype=0是大端，如果是1则是小端
            if [ "$mipstype" = "0" ]; then
                ARCH_TYPE="ppc64"    # 大端
            else
                ARCH_TYPE="ppc64le"  # 小端
            fi
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

ckcmd() {
	command -v sh >/dev/null 2>&1 && command -v $1 >/dev/null 2>&1 || type $1 >/dev/null 2>&1
}

# 检查是否为root用户
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root"
   exit 1
fi

# 检查是不是linux
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

# 检测并配置IP转发
echo "正在检查 IP 转发状态..."
if [ "$(cat /proc/sys/net/ipv4/ip_forward)" = "0" ]; then
    echo "-----------------------------------------------"
    echo -e "\033[33m检测到系统尚未开启IP转发，局域网设备将无法正常连接网络。\033[0m"
    read -p "是否立即开启IP转发？(Y/n) " res
    case "$res" in
        [nN])
            echo "未开启IP转发，部分功能可能无法正常使用"
            ;;
        *)
            echo "正在开启IP转发..."
            if ! grep -q "net.ipv4.ip_forward=1" /etc/sysctl.conf; then
                echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf || {
                    echo "写入 sysctl.conf 失败"
                    exit 1
                }
            fi
            sysctl -w net.ipv4.ip_forward=1 || {
                echo "应用 sysctl 配置失败"
                exit 1
            }
            echo "IP转发已成功开启"
            ;;
    esac
fi

# 检测系统类型: OpenWrt procd 或 systemd
INIT_SYSTEM=""
if [ -f /etc/rc.common ] && [ "$(cat /proc/1/comm)" = "procd" ]; then
    INIT_SYSTEM="procd"
    echo "检测到 OpenWrt procd 系统"
elif ckcmd systemctl && [ "$(cat /proc/1/comm)" = "systemd" ]; then
    INIT_SYSTEM="systemd"
    echo "检测到 systemd 系统"
else
    echo "未能识别系统启动方式，将采用默认配置"
    INIT_SYSTEM="unknown"
fi

# 下载并安装 uaProxy
UaProxy_Name="uaProxy_Linux_${ARCH_TYPE}"
DOWNLOAD_URL="https://github.com/huhu415/uaProxy/releases/latest/download/${UaProxy_Name}.tar.gz"

# 检查压缩包是否存在
if [ ! -f "${UaProxy_Name}.tar.gz" ]; then
    echo "正在下载 uaProxy.tar.gz..."
    # 下载
    wget -q "$DOWNLOAD_URL" || {
        echo "下载失败,请检查网络连接或下载链接是否正确"
        exit 1
    }
else
    echo "压缩包已存在,跳过下载"
fi

# 删除已经存在的解压后的文件夹
rm -rf $UaProxy_Name

# 解压
tar -xzf $UaProxy_Name.tar.gz || {
    echo "解压失败"
    exit 1
}

# 安装主程序
mv $UaProxy_Name/uaProxy /usr/sbin/ || {
    echo "移动 uaProxy 到 /usr/sbin/ 失败"
    exit 1
}
chmod +x /usr/sbin/uaProxy || {
    echo "修改 uaProxy 权限失败"
    exit 1
}

# 根据系统类型安装不同的服务文件
case $INIT_SYSTEM in
    procd)
        mv $UaProxy_Name/shell/uaProxy.procd /etc/init.d/ || {
            echo "移动 uaProxy.procd 到 /etc/init.d/ 失败"
            exit 1
        }
        chmod +x /etc/init.d/uaProxy.procd || {
            echo "修改 uaProxy.procd 权限失败"
            exit 1
        }
        ;;
    systemd)
        mv $UaProxy_Name/shell/uaProxy.service /etc/systemd/system/ || {
            echo "移动 uaProxy.service 到 /etc/systemd/system/ 失败"
            exit 1
        }
        systemctl daemon-reload || {
            echo "重载 systemd 配置失败"
            exit 1
        }
        ;;
    *)
        echo "警告: 未能识别系统类型，跳过服务安装, 自行配置开机自启"
        ;;
esac

rm -rf $UaProxy_Name
rm $UaProxy_Name.tar.gz


# 3. 配置 iptables 规则
echo "正在配置 iptables 规则..."
# iptables -t nat -F // 清空 nat 表
# iptables -t nat -X uaProxy
iptables -t nat -L uaProxy >/dev/null 2>&1 || iptables -t nat -N uaProxy

iptables -t nat -A uaProxy -d 192.168.0.0/16 -j RETURN || {
    echo "添加 iptables 规则失败 (1/5)"
    exit 1
}
iptables -t nat -A uaProxy -p tcp -j RETURN -m mark --mark 0xff || {
    echo "添加 iptables 规则失败 (2/5)"
    exit 1
}
iptables -t nat -A uaProxy -p tcp -j REDIRECT --to-ports 12345 || {
    echo "添加 iptables 规则失败 (3/5)"
    exit 1
}
iptables -t nat -A PREROUTING -p tcp -j uaProxy || {
    echo "添加 iptables 规则失败 (4/5)"
    exit 1
}
iptables -t nat -A OUTPUT -p tcp -j uaProxy || {
    echo "添加 iptables 规则失败 (5/5)"
    exit 1
}

# 配置开机启动
echo "正在配置开机启动..."
case $INIT_SYSTEM in
    procd)
        /etc/init.d/uaProxy.procd enable
        echo "已配置 OpenWrt procd 启动项"
        ;;
    systemd)
        systemctl enable uaProxy.service
        echo "已配置 systemd 启动项"
        ;;
    *)
        echo "未能配置开机启动，请手动配置"
        ;;
esac

# 启动服务
echo "正在启动服务..."
case $INIT_SYSTEM in
    procd)
        /etc/init.d/uaProxy.procd start
        ;;
    systemd)
        systemctl start uaProxy.service
        ;;
    *)
        echo "请手动启动服务"
        ;;
esac

echo "安装完成！"
echo "----------------------------------------"
echo "使用说明："
case $INIT_SYSTEM in
    procd)
        echo "1. 使用 '/etc/init.d/uaProxy.procd {start|stop|restart|status}' 控制服务"
        echo "2. 使用 'logread | grep uaProxy' 查看运行日志"
        ;;
    systemd)
        echo "1. 使用 'systemctl {start|stop|restart} uaProxy' 控制服务"
        echo "2. 使用 'systemctl status uaProxy.service' 查看运行状态"
        ;;
    *)
        echo "1. 请根据您的系统手动控制服务"
        echo "2. 请根据您的系统手动查看运行日志"
        ;;
esac
echo "3. iptables中nat表配置完成"
echo "4. IP 转发已配置完成"
echo "----------------------------------------"
