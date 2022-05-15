本质上是 子序列 dp

```Python
# 求子数组按位或操作可能得到的个数
# 时间复杂度是多少呢?
class Solution(object):
    def subarrayBitwiseORs(self, A):
        res = set()  # 存储所有子集(子序列)，子集按照每个元素结尾来划分
        dp = {0}
        for num in A:
            ndp = {num | y for y in dp} | {num}   # 以 num 结尾的子数组的所有或
            res |= ndp  # 更新到结果
            dp = ndp  # 更新到dp
        return len(res)
```
