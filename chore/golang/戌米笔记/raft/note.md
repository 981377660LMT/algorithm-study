# 【解读共识算法Raft】PDF下载

https://mp.weixin.qq.com/s/4szPOiL98uj39_43DrBnxQ

1. 解读共识算法Raft（1）简介和状态简化
   • 相比于Paxos，Raft最大的特性就是易于理解（Understandable）。为了达到这个目标，
   Raft主要做了两方面的事情：
   • 1. 问题分解：把共识算法分为三个子问题，分别是领导者选举（leader election）、日志复制（log replication）、安全性（safety）
   • 2. 状态简化：对算法做出一些限制，减少状态数量和可能产生的变动。
