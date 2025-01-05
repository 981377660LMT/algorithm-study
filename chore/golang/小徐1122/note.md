https://www.bilibili.com/video/BV1k14y117di
https://www.zhihu.com/people/xu-xian-sheng-80-10

# Golang road map

## 入门：

- 刘丹冰老师视频
- mysql
  - gorm 使用：李文周老师技术博客
  - 原理：索引、事务、锁
    周瑜老师 b 站解说 https://www.bilibili.com/video/BV1kA411u734
    推荐：mysql 技术内幕 innodb 存储引擎第 2 版
- redis
  - redigo sdk 开源客户端
  - 小徐 1122 分布式锁
- mq
  - 基于 redis 实现消息队列
  - kafka 架构原理
- web
  - gin 框架，李文周老师技术博客
  - net/http 原理
  - io 多路复用 epoll 原理
  - gin 框架实现原理
- context

## 进阶：

1. 基础知识进阶
   - gmp
   - gc
   - channel
   - mutex
   - map
   - slice
   - waitGroup
2. 分布式理论
   - 分布式锁：golang、redis 分布式锁
   - 共识算法:raft 理论、raft-etcd 案例
   - 分布式事务: 原理、手写 tcc 框架
   - 微服务框架: go-zero 学习攻略
3. 优秀开源项目
   - 协程池: ants
   - rpc 框架:grpc-go
   - kv 组件: etcd，分布式协调者(类似 zookeeper)
4. 数据结构与算法
   - 跳表
   - 日志结构合并树 lsm tree (rocksdb 中使用)
   - 布隆过滤器
   - 前缀树 trie
   - geohash
   - 一致性哈希
