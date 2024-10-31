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
  - 这种检测方式应该只存在于paper上, 没有用这种方式检测了, 我实验过, 就算是一台机器, 也会有时间不同的情况, 所以这种方式不太靠谱
- [x] UA
  - [本项目](https://github.com/huhu415/uaProxy)
- [x] IP-ID
  - 这种检测方式也只存在于paper上, 因为我也实验过, 就算是一台机器, 也不是单调递增的, 不过就算没有这种检测方式[本项目](https://github.com/huhu415/uaProxy)天生就没有这个问题.
- [ ] DPI (可能暂时无解)
  - 只能用规避gfw的方式, 可能要全局代理, 比如使用`v2ray`. 以后再说
