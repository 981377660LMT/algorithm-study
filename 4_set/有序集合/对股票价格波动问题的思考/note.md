最大/最小值
最近使用的值：需要插入的键有序

1. 大根堆+小根堆做 (两个优先队列，不在 map 里面就出列(延迟删除))
2. Multiset/SortedList/HashHeap (将价格全部加入 sortedList,update 时先删除旧的价格再加入新的价格；dict 记录时间与价格的对应关系)
3. SortedDict/TreeMap 键为价格 值为该价格出现的次数 与 **2** 解法相似 (sorteddict 不如 sortedList)

**模拟**
模拟 SortedList:bisect+切片
模拟 Multiset:SortedList 就是 multiset (全 logn);SortedDict 记录次数就是 SortedList

使用 dict 的好处:
使用 sortedList 的好处:

总结:
`id + dict + sortedList 三件套`
