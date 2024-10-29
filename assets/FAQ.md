## ❓Questions

### 如何选择Assets版本？

- **386**: 32位的x86架构处理器
- **amd64**: 64位的x86架构处理器（也称为x86_64）
- **arm**: 32位的ARM架构处理器
- **arm64**: 64位的ARM架构处理器（也称为AArch64）
- **loong64**: 龙芯处理器的64位架构
- **mips**: 32位的大端序MIPS处理器
- **mipsle**: 32位的小端序MIPS处理器
- **mips64**: 64位的大端序MIPS处理器
- **mips64le**: 64位的小端序MIPS处理器
- **ppc64**: 64位的大端序PowerPC处理器
- **ppc64le**: 64位的小端序PowerPC处理器

> 注意, 所有Assets都是Linux版本的, darwin,windows,freebsd等等没有iptables, 无法使用

### 我是路由器, 选择哪个版本？

先查一下自己的cpu的bits, 是32位还是64位的. _这就可以去掉一半了_

再看看是什么架构, 如果是硬路由, 比如ac2100这种, 一般都是`MIPS`类型的
  - 如果是`MIPS`类型的, _那么就看看是大端序还是小端序的, 可以去问[ChatGpt](https://chatgpt.com)或者[kimi](https://kimi.moonshot.cn/chat/)_
  - 如果不是, 那就是`ARM`或者`x86`的

### 校园网常见检测方式
- [x] TTL
  - ```sh
    iptables -t mangle -A POSTROUTING -j TTL --ttl-set 64 # 修改出口 TTL 为 64
    ```
- [x] 时间戳
  - ```sh
    iptables -t nat -N ntp_force_local
    iptables -t nat -I PREROUTING -p udp --dport 123 -j ntp_force_local
    iptables -t nat -A ntp_force_local -d 0.0.0.0/8 -j RETURN
    iptables -t nat -A ntp_force_local -d 127.0.0.0/8 -j RETURN
    iptables -t nat -A ntp_force_local -d 192.168.0.0/16 -j RETURN
    iptables -t nat -A ntp_force_local -s 192.168.0.0/16 -j DNAT --to-destination 192.168.1.1 # 根据你路由器的地址修改
    # 同时记得修改ntp服务器地址, 可以选ntp.aliyun.com, time1.cloud.tencent.com, time.ustc.edu.cn, cn.pool.ntp.org
    ```
- [x] UA
  - [本项目](https://github.com/huhu415/uaProxy)
- [ ] IP-ID
  - 一般学校好像不会检测这个, 以后再说
- [ ] DPI (可能暂时无解)
  - 可能要全局代理, 以后再说
