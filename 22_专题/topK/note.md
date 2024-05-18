查询 topK

- 快速选择 nthElement (O(n)) `一般用于快速求出中位数`
- 堆/SortedList 名次树 (O(nlogk)，数据流，海量数据多路归并)
- 计数排序(O(n),需要已知数据最大最小值)

TODO：数据库翻页查询、堆实现的 lazy-sort、渐进式排序
https://help.aliyun.com/zh/polardb/polardb-for-mysql/user-guide/implementation-of-topk-operator-in-column-store-index?spm=a2c4g.11186623.0.0.4b0735d7BcYE8w
