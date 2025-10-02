剩余 21 篇 TODO

- etcd watch
  https://mp.weixin.qq.com/s/2TEgbOoX36PwSWzbKq0Qsg
  ~~etcd watch 机制源码解析——客户端篇~~

  https://mp.weixin.qq.com/s/-Vxu7jQZ-7ID-4oUF_0Agg
  ~~etcd watch 机制源码解析——服务端篇~~

- sql
  https://mp.weixin.qq.com/s/ojDRfrotU8ByOTIYFZxF0g
  ~~Golang sql 标准库源码解析~~

  https://mp.weixin.qq.com/s/V0HFl0NJFbGS_InTURzMLQ
  ~~Golang mysql 驱动源码解析~~

  https://mp.weixin.qq.com/s/plzG1mCK8yZwVQOSKZi2XQ
  ~~gorm 框架使用教程~~

  https://mp.weixin.qq.com/s/STFnyke1NX8Ag8COlHwaLA
  ~~gorm 框架原理&源码解析~~

- mq
  https://mp.weixin.qq.com/s/a_o7cUxhb9XC_fNecL7WAA
  ~~万字解析 go 语言分布式消息队列 nsq~~
- etcd bolt
  https://mp.weixin.qq.com/s/oL_G8H_ROSF3TjtzBOGCow
  etcd 存储引擎之主干框架

  https://mp.weixin.qq.com/s/nFlcRJagr-UG6LhXmsp4eA
  etcd 存储引擎之存储设计

  https://mp.weixin.qq.com/s/lqFkUIiabcRb2LAEXdXsSA
  etcd 存储引擎之 b+树实现

  https://mp.weixin.qq.com/s/EB-sQXxHtaqelneJ6sUHng
  etcd 存储引擎之事务实现

- redis
  https://mp.weixin.qq.com/s/tKtmhCNtc696a87NBbo1Dw
  基于 go 实现 redis 之主干框架

  https://mp.weixin.qq.com/s/KuRNKJXUtDdWlAP8eyKtsA
  基于 go 实现 redis 之指令分发

  https://mp.weixin.qq.com/s/Yx1R2Rai34W59vTWSS6pFQ
  基于 go 实现 redis 之存储引擎

  https://mp.weixin.qq.com/s/tOyZQ2UPdRrBRlQDCKgWIQ
  基于 go 实现 redis 之数据持久化

- 一致性缓存
  https://mp.weixin.qq.com/s/h1oi92BbdFdTGtey0wQLLQ
  一致性缓存理论分析与技术实战

- innodb
  https://mp.weixin.qq.com/s/VTK03xFSfTlFxECC_xhJcw
  万字解析 mysql innodb 锁机制实现原理

  https://mp.weixin.qq.com/s/dEgTmNioqjx5IyDTX4KE4g
  万字解析 mysql innodb 事务实现原理

- 协程
  https://mp.weixin.qq.com/s/cWBRZ9Zha8qN_A9C_epn1A
  C++ 从零实现协程调度框架

- netpoll
  https://mp.weixin.qq.com/s/_FTvpvLIWfYzgNhOJgKypA
  万字解析 golang netpoll 底层原理

- 对象池
  https://mp.weixin.qq.com/s/T2ui02f05l6prn7jYxn2RQ
  go/c++ 万字解析对象池技术原理与源码实战

- sync.Map
  https://mp.weixin.qq.com/s/YU_nACyZEkQRIHXWO6He5A
  go/c++ 万字解析 sync.Map 技术原理与源码实战
- gmp
  https://mp.weixin.qq.com/s/BR6SO7bQF4UXQoRdEjorAg
  温故知新——Golang GMP 万字洗髓经
- cloudwego
  https://mp.weixin.qq.com/s/BUPitmZR1hxXIPlFfAUWgw
  万字学习笔记：cloudwego/netpoll

  https://mp.weixin.qq.com/s/pEUMJccbOmDn65xVUancGg
  万字学习笔记：cloudwego/kitex

  https://mp.weixin.qq.com/s/HxTDCcT79FZuDoMWXK-Fow
  cloudwego/hertz 原理浅析

- misc
  https://mp.weixin.qq.com/s/yXv_kfgWFRT5MyWUMyeOuQ
  ~~初窥门径——前端 React 项目实战学习笔记~~

  https://mp.weixin.qq.com/s/WyXIUjAUeOAVg0uhwj1WIg
  浅谈 DDD 领域驱动设计架构

  https://mp.weixin.qq.com/s/Qwmo9FHCF010-0rMUbyuww
  Golang pprof 案例实战与原理解析

---

history

# Go 文章分类汇总

## 1. Go 语言核心原理

### 并发模型与内存管理

- [Golang GMP 原理](https://mp.weixin.qq.com/s/jIWe3nMP6yiuXeBQgmePDg)
- [Golang 内存模型与分配机制](https://mp.weixin.qq.com/s/2TBwpQT5-zU4Gy7-i0LZmQ)
- [Golang 垃圾回收原理分析](https://mp.weixin.qq.com/s/TdekaMjlf_kk_ReyPvoXiQ)
- [Golang 垃圾回收源码走读](https://mp.weixin.qq.com/s/Db19tKNer8D6FX6UG-Yujw)

### 数据结构实现

- [Golang map 实现原理](https://mp.weixin.qq.com/s/PT1zpv3bvJiIJweN3mvX7g)
- [Golang sync.Map 实现原理](https://mp.weixin.qq.com/s/nMuCMA8ONnhs1lsTVMcNgA)
- [你真的了解 go 语言中的切片吗？](https://mp.weixin.qq.com/s/uNajVcWr4mZpof1eNemfmQ)
- [Golang Channel 实现原理](https://mp.weixin.qq.com/s/QgNndPgN1kqxWh-ijSofkw)

### 并发控制

- [Golang context 实现原理](https://mp.weixin.qq.com/s/AavRL-xezwsiQLQ1OpLKmA)
- [Golang 单机锁实现原理](https://mp.weixin.qq.com/s/5o0pR0RDaasKh4veXTctVg)
- [Go 并发编程之 sync.WaitGroup](https://mp.weixin.qq.com/s/oPVmOT3rpkulraz_pDqWsA)

### 网络与 HTTP

- [解析 Golang 网络 IO 模型之 EPOLL](https://mp.weixin.qq.com/s/xt0Elppc_OaDFnTI_tW3hg)
- [Golang HTTP 标准库实现原理](https://mp.weixin.qq.com/s/zFG6_o0IKjXh4RxKmPTt4g)

## 2. 设计模式与架构

### 设计模式

- [Golang 设计模式之单例模式](https://mp.weixin.qq.com/s/KRgNwJt1C7q2ckeqCu9pCQ)
- [Golang 设计模式之工厂模式](https://mp.weixin.qq.com/s/Fp0KxoXc4y3z8cgc7IzD-Q)
- [Golang 设计模式之观察者模式](https://mp.weixin.qq.com/s/QOXh86eX8z5Ts4O1pky44g)
- [Golang 设计模式之装饰器模式](https://mp.weixin.qq.com/s/XG2G1O67o-p_u_MPj5N0oQ)
- [Golang 设计模式之适配器模式](https://mp.weixin.qq.com/s/g1xKCmRMqxP9DmtW208A6w)
- [Golang 设计模式之建造者模式](https://mp.weixin.qq.com/s/PuDyI0hA2xZsCwh7imyxoQ)
- [Golang 设计模式之责任链模式](https://mp.weixin.qq.com/s/JCxPUg1MHz7ITCxKyli-mw)

### 依赖注入

- [低配 Spring—Golang IOC 框架 dig 原理解析](https://mp.weixin.qq.com/s/bireIkWWTQUdgc-UJEhN5g)

## 3. 数据结构与算法

### 高性能数据结构

- [基于 golang 从零到一实现跳表](https://mp.weixin.qq.com/s/fvfz6bdvsZJtGsdL0MPYoA)
- [如何实现一个并发安全的跳表](https://mp.weixin.qq.com/s/7VhioGP007LDQnZ_w8GBBQ)
- [基于 Golang 实现前缀树 Trie](https://mp.weixin.qq.com/s/_4K-zDZgCPvSBmjHbj6GGA)
- [布隆过滤器技术原理及应用实战](https://mp.weixin.qq.com/s/_dtmItfAnHn6x8s0zSzFLA)

### LSM 树实现

- [初探 rocksdb 之 lsm tree](https://mp.weixin.qq.com/s/kqpBZ2aCC0CGvvL2Lm6mzA)
- [基于 go 实现 lsm tree 之主干框架](https://mp.weixin.qq.com/s/KpkiBQNycLoDskUr00WeEg)
- [基于 go 实现 lsm tree 之 memtable 结构](https://mp.weixin.qq.com/s/WUm9mu8XtMgb7iJDsYYDNg)
- [基于 go 实现 lsm tree 之 sstable 结构](https://mp.weixin.qq.com/s/-1TuGjDXNj7z1ZSZsfQjYg)
- [基于 go 实现 lsm tree 之 level sorted merge 流程](https://mp.weixin.qq.com/s/dwuYGRRxJgEhfHkRyzyPjQ)

### 分布式算法

- [一致性哈希算法原理解析](https://mp.weixin.qq.com/s/NZNVCZF8jiiPTAQ8Ievyrg)
- [从零到一落地实现一致性哈希算法](https://mp.weixin.qq.com/s/e5Q4pYmfqpa9Ix-GIKpvMA)
- [GeoHash 技术原理及应用实战](https://mp.weixin.qq.com/s/5AwyHocmwqVQOsNk-Al0xA)

### 定时器与时间轮

- [基于协程池架构实现分布式定时器 XTimer](https://mp.weixin.qq.com/s/gfiAm4NrcY_PaRNrQ1P2vw)
- [Golang 协程池 Ants 实现原理](https://mp.weixin.qq.com/s/Uctu_uKHk5oY0EtSZGUvsA)
- [基于 golang 从零到一实现时间轮算法](https://mp.weixin.qq.com/s/0HwuYTTe9cdT46advEZe0w)

## 4. 分布式系统

### 分布式共识

- [两万字长文解析 raft 算法原理](https://mp.weixin.qq.com/s/nvg9J4ky9mz-dFVi5CyYWg)
- [raft 工程化案例之 etcd 源码实现](https://mp.weixin.qq.com/s/jsJ3_E_5IOs4_rPDM5axzQ)

### 分布式锁

- [Golang 分布式锁技术攻略](https://mp.weixin.qq.com/s/KYiZvFRX0CddJVCwyfkLfQ)
- [redis 分布式锁进阶篇](https://mp.weixin.qq.com/s/3zuATaua6avMuGPjYEDUdQ)

### 分布式事务

- [万字长文漫谈分布式事务实现原理](https://mp.weixin.qq.com/s/Z-ZY9VYUzNER8iwk80XSxA)
- [从零到一搭建 TCC 分布式事务框架](https://mp.weixin.qq.com/s/aTQ6mgVUbUqn69NLmXh_7A)

## 5. 微服务与框架

### gRPC

- [grpc-go 服务端使用介绍及源码分析](https://mp.weixin.qq.com/s/OiQ5I1TLex3G-AVBzNN51Q)
- [grpc-go 客户端源码走读](https://mp.weixin.qq.com/s/IkYXT39p1xvwg1AIWbdX3w)
- [基于 etcd 实现 grpc 服务注册与发现](https://mp.weixin.qq.com/s/x-vC1gz7-x6ELjU-VYOTmA)

### Web 框架

- [解析 Gin 框架底层原理](https://mp.weixin.qq.com/s/x8i9HvAzIHNbHCryLw5icg)

## 6. 中间件与存储

### 消息队列

- [万字长文解析如何基于 Redis 实现消息队列](https://mp.weixin.qq.com/s/MSmipbE5cyK2_m5iiKv7pw)

### 版本控制

- [万字串讲 git 版本控制底层原理及实战分享](https://mp.weixin.qq.com/s/R73aUyUHJy_j4tMnq4qXGw)
