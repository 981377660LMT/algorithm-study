## Lecture 1: 关系模型

## Lecture 2: Modern SQL 概述

---

Query Planning
Operator Execution
Access Methods
Buffer Pool Manager
Disk Manager

## Lecture 3: 存储引擎1

## Lecture 4: 存储引擎2

## Lecture 5: 数据库存储 II – 日志结构存储与磁盘优化

## Lecture 6: 存储模型与数据压缩

## Lecture 7: 哈希表与散列索引

## Lecture 8: B+树索引

## Lecture 9: 索引并发

## Lecture 10: 排序与聚合算法

没学“外部排序”，无法处理超出内存的大规模ORDER BY；
没学聚合算法，就不知道数据库怎么做GROUP BY 才又快又省资源。

## Lecture 11: 连表算法 - 排序 vs 哈希 vs 嵌套循环

三大经典连接算法——嵌套循环、排序-合并和哈希连接
“没有索引先排序，等值连接用哈希，深度优化拿索引。”

## Lecture 12: 查询执行 I – 逻辑计划与执行器架构

## Lecture 13: 查询执行 II – 高级执行器与优化

## Lecture 14: 查询优化

---

## Lecture 15: 事务与并发控制理论

## Lecture 16: 两阶段锁并发控制 (Strict 2PL)

## Lecture 17: 基于时间戳的并发控制 (OCC & MVCC)

## Lecture 18: 多版本并发控制 (MVCC) 实现

## Lecture 19: 数据库日志

## Lecture 20: 数据库恢复 – 日志与重做/撤销

## Lecture 21: 数据库恢复 – 崩溃恢复过程

## Lecture 22: 分布式 OLTP 系统案例

## Lecture 23: 分布式 OLAP 系统案例
