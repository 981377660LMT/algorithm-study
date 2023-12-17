# 统计区间内元素个数的两种方法

1. 前缀和统计区间内元素个数 -> alphaPresum ，`O(26*n)预处理，O(1)查询`
2. 哈希表+二分统计区间内元素个数 -> RangeFreqQuery，`O(n) 预处理，O(logn)查询`
