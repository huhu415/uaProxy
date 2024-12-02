#!/bin/sh

SERVICE_NAME="uaProxy.procd"
INIT_SCRIPT="/etc/init.d/$SERVICE_NAME"

echo "正在配置 iptables 规则..."
iptables -t nat -F
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

# 确保服务脚本可执行
if [ -f "$INIT_SCRIPT" ]; then
    chmod +x "$INIT_SCRIPT"
fi

# 启用服务
if [ -x "$INIT_SCRIPT" ]; then
    "$INIT_SCRIPT" enable
    "$INIT_SCRIPT" start
fi

exit 0
