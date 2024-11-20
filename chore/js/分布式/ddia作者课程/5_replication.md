## Replication

A node that has a copy of the data is called a `replica`.

- Idempotence
  ![alt text](image-15.png)

  - 性质：idempotent requests `can be retried without deduplication`
  - 三种retry语义(semantics)：

    - at-most-once
      send request, `dont retry`, update may not happen
    - at-least-once
      `retry request until ack`, may repeat update
    - exactly-once
      `retry + idempotence` or `retry + dedup`

  - Adding and then removing again problem
    ![alt text](image-16.png)
    因为：`虽然 f(f(x)) = f(x)，但是 f(g(f(x)))!=g(f(x))`
  - Another problem with adding and removing
    ![alt text](image-17.png)
    最终状态相同，但是用户意图不同

    解决方案：
    **timestamps + tombstones**
    ![alt text](image-18.png)

    `每次操作加上版本(eg:时间戳)，删除使用标记删除(tombstones)`

- Reconciling replicas (协调副本)
  ![alt text](image-19.png)
  根据版本更新副本值

- Concurrent writes by different clients
  ![alt text](image-20.png)
  两种解决方式：

  1. Last write wins (LWW)
     如果使用全序时钟，可以保留最后一个。
     但问题是：data loss。
  2. Multi-value register
     如果使用偏序时钟，根据并发关系/偏序关系决定保留后一个还是两个都保留。

## Quorums

## State Machine Replication

- 串行事务
- 区块链、智能合约
