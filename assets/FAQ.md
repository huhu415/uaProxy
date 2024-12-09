## ❓常见问题

### 如何选择 Assets 版本？

1. **确认 CPU 架构和位数**：使用命令 `uname -m` 查看是arm64, arm或者mips, mips64等
2. **确认 CPU 大端小段**：一般都是大端, 很少很少有小端, 除非是单片机这种, 可咨询 [ChatGPT](https://chatgpt.com) 或 [kimi](https://kimi.moonshot.cn/chat/)。

处理器架构对应版本：

| 版本                 | 处理器架构        |
| -------------------- | ----------------- |
| **386**              | 32 位 x86         |
| **amd64(x86_64)**    | 64 位 x86         |
| **arm**              | 32 位 ARM         |
| **arm64（AArch64）** | 64 位 ARM         |
| **loong64**          | 64 位龙芯         |
| **mips**             | 32 位大端 MIPS    |
| **mipsle**           | 32 位小端 MIPS    |
| **mips64**           | 64 位大端 MIPS    |
| **mips64le**         | 64 位小端 MIPS    |
| **ppc64**            | 64 位大端 PowerPC |
| **ppc64le**          | 64 位小端 PowerPC |

> 注意：所有 Assets 均为 Linux 版本，因 darwin、windows、freebsd 等系统没有 iptables，无法使用。

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
  - 只能用规避[gfw](https://zh.wikipedia.org/wiki/%E9%98%B2%E7%81%AB%E9%95%BF%E5%9F%8E)的方式, 可能要全局代理, 比如使用`v2ray`. 以后再说
