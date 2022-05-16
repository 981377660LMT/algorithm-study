本质上是 dp

```Python
# 求子数组按位或操作可能得到的个数
# ndp 的 元素个数有限时 可以这样做
# O(nlogA)
class Solution(object):
    def subarrayBitwiseORs(self, A):
        res = set()
        dp = {0}
        for num in A:
            ndp = {num | y for y in dp} | {num}   # 以 num 结尾的子数组的所有或
            dp = ndp  # 滚动替换
            res |= ndp
        return len(res)
```

```Python
# 求子序列(子集)按位或操作可能得到的个数

# O(2^n)
class Solution(object):
    def subarrayBitwiseORs(self, A):
        dp = {0}
        for num in A:
            ndp = {num | y for y in dp}   # 以 num 结尾的子集(序列)的所有或
            dp |= ndp  # 直接加到dp，形成新的子集
        return len(dp)
```
