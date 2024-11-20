什么是分布式系统？分布式系统会遇到哪些问题？既然有这些问题为什么还要有分布式系统

## 什么是分布式系统

A system running on several nodes connected by a network, characterized by `partial failure`.
通过网络连接的多个节点上运行的系统，其特征是`部分故障`。

## Failures 故障

两种处理哲学

1. High-Performance Computing (HPC) philosophy
   高性能计算 (HPC) 理念

   - treat partial failure as total failure
     将部分失败视为完全失败
   - checkpointing
     检查点

2. “cloud computing” philosophy
   “云计算”理念
   - treat partial failure as expected and work around it
     按预期对待部分失败
   - “everything fails, all the time”
     “一切都会失败，一直如此”

## 故障场景

考虑一个由两台机器 M1 和 M2 组成的简单系统。假设 M1 向 M2 询问某个变量x的值。

```
M1 -- x? --> M2
   <-- 5 ---
```

一些可能的失败是：

- request from M1 is lost
  M1 的请求丢失
- request from M1 is slow
  M1 的请求很慢
- M2 crashes
  M2 崩溃
- M2 is slow to respond
  M2反应慢
- response from M2 is slow or lost
  M2 的响应缓慢或丢失
- Byzantine faults response/etc
  corrupt transmission/malicious
  拜占庭错误(系统中的某些节点可能会表现出任意的、不可预测的行为)
  传输损坏/恶意响应等

这些问题的关键：
**如果您向另一个节点发送请求但没有收到响应，是不可能知道原因的。**
真实的系统尝试使用`超时`来解决这个问题。

## Timeouts 超时

After a certain amount of time, assume failure and try again.
在一段时间后，假设失败并重试.
要求操作必须是`幂等`的，即重试不会产生副作用。
这不是一个完美的解决方案.

一般选择的超时时间是(2d+r), 其中d是消息传输的延迟，r是消息处理的时间。
然而，在现实生活中通常无法得到d和r的准确值，只能得到统计数据。

Distributed systems have to account for both partial failure and unbounded latency.
`因此，分布式系统必须考虑部分故障和无限延迟。`

## 为什么使用分布式系统

Why do we make systems distributed and deal with all of this?

- inherently distributed
  本质上是分布式的，例如发送短信
- make things faster (better perf)
  更快解决问题
- more data than can fit on one machine (bigger problems)
  无法在一台机器上容纳更多数据
- reliability (more copies of data)
  可靠性(更多数据`副本`)
- throughput (data is physically closer)
  吞吐量

## 为什么不使用分布式系统

- communication may fail
  通信可能失败，甚至不知道失败了
- processes may fail
  进程可能失败，甚至不知道失败了
- non-deterministic
  非确定性
