# uaProxy

利用clash, v2ray的方案(iptables-redir)实现对所有流量的监控,
从而找出http流量后, 修改其中的`User-Agent`字段, 从而实现对所有http流量的`User-Agent`修改.

> 注意, 因为是利用了iptables的REDIRECT功能, 所以不能和clash, v2ray等软件同时使用, 会有冲突.

> 但这样做也更纯净, 性能最快, 我觉得应该是这个需求的最佳实现方案了.
