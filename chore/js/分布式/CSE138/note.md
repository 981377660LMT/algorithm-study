https://cse138-notes.readthedocs.io/en/latest/
https://www.bilibili.com/video/BV1AU4y1h76B

# 分布式系统学习笔记 (CSE 138)

## 目录

### 基础概念

- [Introduction](#introduction)
- [Failures](#failures)
- [Timeouts](#timeouts)
- [Why?](#why)

### 时间与同步

- [Time](#time)
- [Logical Clocks](#logical-clocks)
- [Lamport Diagrams](#lamport-diagrams)
- [Lamport Clocks](#lamport-clocks)
- [Vector Clocks](#vector-clocks)

### 网络与通信

- [Network Models](#network-models)
- [State and Events](#state-and-events)
- [Partial Order](#partial-order)
- [Protocol](#protocol)
- [FIFO Delivery](#fifo-delivery)
- [Causal Delivery](#causal-delivery)
- [Totally Ordered Delivery](#totally-ordered-delivery)

### 分布式系统快照

- [Snapshots of Distributed Systems](#snapshots-of-distributed-systems)
- [Chandy-Lamport Algorithm](#chandy-lamport-algorithm)

### 容错与可靠性

- [Safety & Liveness](#safety--liveness)
- [Fault Models](#fault-models)
- [Two Generals Problem](#two-generals-problem)
- [Fault Tolerance](#fault-tolerance)
- [Reliable Delivery, Take 2](#reliable-delivery-take-2)
- [Reliable Broadcast](#reliable-broadcast)

### 复制与一致性

- [Replication](#replication)
- [Primary-Backup Replication](#primary-backup-replication)
- [Chain Replication](#chain-replication)
- [Total Order v. Determinism](#total-order-v-determinism)
- [Consistency](#consistency)
- [Coordination](#coordination)
- [Active v. Passive Replication](#active-v-passive-replication)

### 共识算法

- [Consensus](#consensus)
  - [Properties](#properties)
  - [Paxos](#paxos)
  - [Multi-Paxos](#multi-paxos)
  - [Paxos: Fault Tolerance](#paxos-fault-tolerance)
  - [Other Consensus Protocols](#other-consensus-protocols)

### 一致性模型

- [Consistency](#consistency)
- [Eventual Consistency](#eventual-consistency)
- [Quorum Consistency](#quorum-consistency)
- [Dynamo](#dynamo)

### 分布式系统设计

- [Sharding](#sharding)
- [Consistent Hashing](#consistent-hashing)
- [Heterogeneous Distributed Systems](#heterogeneous-distributed-systems)
- [First-Order Distributed Systems](#first-order-distributed-systems)
- [Monolith](#monolith)

### MapReduce

- [MapReduce](#mapreduce)
- [Word Count](#word-count)
- [What Could Go Wrong](#what-could-go-wrong)
- [MapReduce @ Google](#mapreduce--google)

### 理论与数学

- [The Cost of Consensus](#the-cost-of-consensus)
- [Strong Convergence](#strong-convergence)
- [Upper Bounds](#upper-bounds)
- [Examples](#examples)
