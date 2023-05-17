# 6.851: Advanced Data Structures

## Lecture 01

四种持久化等级

## Lecture 02

三种追溯化等级

## Lecture 03

点位置和范围搜索
确定用户单击了哪个 GUI 元素、一组 GPS 坐标位于哪个城市以及某些类型的数据库查询等应用程序。
三维情况下可以预处理后每次查询 O(logn)

## Lecture 04

O（log n）通过分数级联进行三维正交范围搜索

## Lecture 05

动态最优猜想是数据结构中最古老和最大的开放问题之一，它提出了一个基本问题：是否有一个最好的二叉搜索树？
splay tree 符合

## Lecture 06

Tango trees O(loglogn)

## Lecture 07

缓存高效的数据结构
这里的一个经典结果是，B 树善于利用数据在缓存和主内存之间以及主内存和磁盘之间以块的形式传输，等等。B 树实现了 N 个项目的 O（log B N） 插入/删除/前置/后继器和大小为 B 的内存块传输。最近和令人惊讶的是，即使您不知道 B 是什么，也可以实现相同的性能，或者换句话说，对于具有所有 B 值的所有架构，也可以同时实现相同的性能。这是“忽略缓存(cache-oblivious)”模型的结果.

## Lecture 08

缓存的数据结构。

## Lecture 09

## Lecture 10

不同的哈希函数

## Lecture 11

前驱后继快速查找

## Lecture 12

融合树(fusion tree)

## Lecture 13

## Lecture 14

## Lecture 15

## Lecture 16

## Lecture 17

succient DS 使用空间很少
在 2n + o（n） 位中存储 n 节点二进制 trie
这种 trie 是基础

## Lecture 18

## Lecture 19

## Lecture 20

动态连通性

## Lecture 21

## Lecture 22
