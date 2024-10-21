# uaProxy

利用clash, v2ray的方案(iptables-redir)实现对所有流量的监控,
找出http流量后, 修改其中的`User-Agent`字段, 从而实现对所有http流量的`User-Agent`修改.
![uaProxy](assets/uaProxy.png)

## 代码思路
流量分为http和非http两种, 非http流量直接转发, http流量则可以通过看前几个字节判断是否是http流量.

如果是http流量, 则利用go官方的webServer的实现方法, 循环使用`http.ReadRequest()`读取http请求, 从中找到`User-Agent`字段, 修改后再写回.

## 使用方法
1. 网关设备开启 IP 转发。
在 `/etc/sysctl.conf` 文件添加一行 `net.ipv4.ip_forward=1` ，执行下列命令生效：`sysctl -p`

2. 运行uaProxy (所有linux都可以用, 这里详细讲openWrt)
  - 脚本
    1. 执行``
  - 手动
    - 下载[相应](https://github.com/huhu415/uaProxy/releases)的压缩包, 解压后
      - 把可执行程序放到`/usr/sbin`目录里面
      - 把[脚本文件](assets/uaProxy-openwrt)放到`/etc/init.d`目录里面
    - 执行`chmod +x /etc/init.d/uaProxy-openwrt`, 赋予执行权限
      - 执行`/etc/init.d/uaProxy-openwrt enable`, 开机自启
      - 执行`/etc/init.d/uaProxy-openwrt start`, 启动服务
      - (可选)执行`logread | grep uaProxy` 查看日志; 同时也可以登陆web页面, 在`状态-系统日志`里面看

3. 为了实现所有TCP流量会经过uaProxy, iptables要这样设置
```sh
iptables -t nat -N uaProxy # 新建一个名为 uaProxy 的链
iptables -t nat -A uaProxy -d 192.168.0.0/16 -j RETURN # 直连 192.168.0.0/16
iptables -t nat -A uaProxy -p tcp -j RETURN -m mark --mark 0xff
# 直连 SO_MARK 为 0xff 的流量(0xff 是 16 进制数，数值上等同与上面配置的 255)，此规则目的是避免代理本机(网关)流量出现回环问题
iptables -t nat -A uaProxy -p tcp -j REDIRECT --to-ports 12345 # 其余流量转发到 12345 端口（即 uaProxy默认开启的redir-port）
iptables -t nat -A PREROUTING -p tcp -j uaProxy # 对局域网其他设备进行透明代理
iptables -t nat -A OUTPUT -p tcp -j uaProxy # 对本机进行透明代理, 可以不加
```

### 参数说明:
`--stats` 开启统计信息
- 不开启(默认): 修改所有`http`流量的`UA`为统一字段.
- 开启: 在可执行程序同目录下生成一个`stats-config.csv`文件, 里面记录了不同`User-Agent`字段的访问次数.
  - 如果记录项有`**uaProxy**`前缀, 代表已经检测到特征, 会被修改为统一的UA字段; 否则不会修改.
  - 建议只有在普通模式有问题时, 再开启统计模式, 以免影响性能和反检测效果.


> ‼️注意, 因为是利用了iptables的REDIRECT功能, 所以不能和clash, v2ray等软件同时使用, 会有冲突.
> 但这样做也更纯净, 性能最快, 我觉得应该是这个需求的最佳实现方案了.

[FAQ](assets/FAQ.md)
