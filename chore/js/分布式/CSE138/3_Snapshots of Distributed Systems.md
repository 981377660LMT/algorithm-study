# Snapshots of Distributed Systems 分布式系统快照

Processes have individual state, which is pretty much all the events that have happened up to a given point.
进程具有独立的状态，这几乎是到给定点发生的所有事件。

How do we get the state of an entire distributed system (a “global snapshot”) at a given time?
我们如何获取给定时间整个分布式系统的状态（“全局快照”）？

We can’t use time-of-day clocks, since they aren’t guaranteed to be synchronized across all processes.
我们不能使用time-of-day时钟，因为它们不能保证在所有进程之间同步。

Property that we want: If we have events A and B where
, and B is in the snapshot, then A is also in the snapshot.
`我们想要的属性：我们有事件A和B ，如果B在快照中，那么A也在快照中。`
一些快照示例：
![alt text](image-18.png)

## Chandy-Lamport Algorithm

Chandy-Lamport Algorithm 是一种用于`在分布式系统中获取一致性快照（Snapshot）的算法`。该算法由 K. Mani Chandy 和 Leslie Lamport 于 1985 年提出，旨在在不停止系统运行的情况下，获取系统的全局状态。

### channel 通道

connection from one process to another, with FIFO ordering
从一个进程到另一个进程的连接，采用 FIFO 排序
