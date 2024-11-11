#!/bin/sh

# 检查是否为root用户
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root"
   exit 1
fi

# 停止服务
/etc/init.d/uaProxy-openwrt stop

# 清理 iptables 规则
echo "正在清理 iptables 规则..."
iptables -t nat -F  # 清空 nat 表
iptables -t nat -X uaProxy || {
    echo "删除 iptables 链规则失败"
    # 这里不退出，继续执行后续卸载步骤
}

# 删除可执行文件
rm -f /usr/sbin/uaProxy || {
    echo "删除 uaProxy 失败"
    exit 1
}

# 删除启动脚本
rm -f /etc/init.d/uaProxy-openwrt || {
    echo "删除 uaProxy-openwrt 失败"
    exit 1
}

echo "卸载完成"
