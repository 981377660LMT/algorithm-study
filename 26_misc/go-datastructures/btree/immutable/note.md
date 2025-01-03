Immutable B Tree 不可变 B 树
A btree based on two principles, immutability and concurrency. Somewhat slow for single value lookups and puts, it is very fast for bulk operations. A persister can be injected to make this index persistent.
基于两个原则的 btree，不可变性和并发性。对于单值查找和放置操作来说速度较慢，但对于批量操作非常快。可以注入一个持久化器使该索引持久化。
