#!/bin/sh

remove() {
    printf "\033[32m 执行 uaProxy 包的卸载操作...\033[0m\n"

    # 停止服务
    /etc/init.d/uaProxy.procd stop

    # 清理 uaProxy 相关的 iptables 规则
    echo "正在清理 uaProxy 相关的 iptables 规则..."
    iptables -t nat -F uaProxy   # 清空 uaProxy 链中的规则
    iptables -t nat -X uaProxy   # 删除 uaProxy 链

    echo "卸载完成" > /tmp/postremove-proof
}

# purge() {
#     printf "\033[32m Post Remove purge, deb only\033[0m\n"
#     echo "Purge" > /tmp/postremove-proof
# }

upgrade() {
    printf "\033[32m Post Remove of an upgrade\033[0m\n"
    echo "Upgrade" > /tmp/postremove-proof
}

echo "$@"

action="$1"

case "$action" in
  "0" | "remove")
    remove
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    remove
    # purge
    ;;
  *)
    printf "\033[32m Alpine\033[0m"
    remove
    ;;
esac
