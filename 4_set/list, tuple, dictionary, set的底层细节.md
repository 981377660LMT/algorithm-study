<!-- https://blog.csdn.net/siyue0211/article/details/80560783 -->

在 CPython 中，列表被实现为**长度可变的数组**。

CPython 使用伪随机探测(pseudo-random probing)的散列表(hash table)作为字典的底层数据结构
由于这个实现细节，只有可哈希的对象才能作为字典的键。

CPython 中集合和字典非常相似。事实上，集合被实现为带有空值的字典，只有键才是实际的集合元素。此外，集合还利用这种没有值的映射做了其它的优化。
