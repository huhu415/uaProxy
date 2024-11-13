#!/bin/sh

ckcmd() {
	command -v sh >/dev/null 2>&1 && command -v $1 >/dev/null 2>&1 || type $1 >/dev/null 2>&1
}

# 检查是否为root用户
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root"
   exit 1
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

# 根据系统类型停止服务并删除启动脚本
case $INIT_SYSTEM in
    "procd")
        # OpenWrt系统
        /etc/init.d/uaProxy stop
        rm -f /etc/init.d/uaProxy.procd
        ;;
    "systemd")
        # systemd系统
        systemctl stop uaProxy
        systemctl disable uaProxy
        rm -f /etc/systemd/system/uaProxy.service
        systemctl daemon-reload
        ;;
    *)
        echo "未知系统类型，跳过服务停止和启动脚本删除"
        ;;
esac

# 清理 iptables 规则
echo "正在清理 iptables 规则..."
iptables -t nat -F  # 清空 nat 表
iptables -t nat -X uaProxy # 删除自定义链

# 删除可执行文件
rm -f /usr/sbin/uaProxy || {
    echo "删除 uaProxy 失败"
    exit 1
}


echo "卸载完成"
