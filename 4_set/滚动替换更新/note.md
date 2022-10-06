**本质上是 相邻行间 dp 用滚动数组优化(想象 bfs 图)**

```Python
# 求子数组按位或操作可能得到的个数
# ndp 的 元素个数有限时 可以这样做
# O(nlogA)
class Solution(object):
    def subarrayBitwiseORs(self, A):
        res = set()
        dp = set()
        for cur in A:
            ndp = {cur | pre for pre in dp|{0}}   # 以 cur 结尾的子数组的所有或
            dp = ndp  # 滚动替换
            res |= ndp
        return len(res)
```

```Python
# 求子序列(子集)按位或操作可能得到的个数

# O(2^n)
class Solution(object):
    def subarrayBitwiseORs(self, A):
        dp = set([0]) # 包含空集
        for cur in A:
            ndp = {cur | pre for pre in dp|{0}}   # 以 cur 结尾的子集(序列)的所有或
            dp |= ndp  # 直接加到dp，形成新的子集
        return len(dp)
```

**当 num 很小时 set 可以用位运算代替**
[6000ms 优化到 60ms](LCP%2065.%20%E8%88%92%E9%80%82%E7%9A%84%E6%B9%BF%E5%BA%A6.py)

```Python
def check(mid: int) -> bool:
  """
  给数组元素添加正负号后,max(preSum) - min(preSum) <= mid 是否成立
  """
  mask = (1 << (mid + 1)) - 1
  dp = mask  # 枚举0-mid起点
  for num in nums:
      dp = ((dp << num) | (dp >> num)) & mask  # & mask 去除超出边界的非法状态
  return dp != 0
```

golang 里用 value 为空结构体的 map 来模拟 set
**当 value 为一个布尔值时,可以用 set 代替 dp 数组**
