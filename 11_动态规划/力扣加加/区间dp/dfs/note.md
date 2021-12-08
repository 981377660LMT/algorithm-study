1. **左右端点为参数型**
   `375. 猜数字大小 II copy.py`
   `730. 统计不同回文子序列-经典模板.py`
   `1000. 合并石头的最低成本.py `
   `1246. 删除回文子数组.py`
   `1547. 切棍子的最小成本.py`
   `2019. 解出数学表达式的学生分数.py`

```Python
class Solution:
    def minimumMoves(self, arr: List[int]) -> int:
        @lru_cache(maxsize=None)
        def dfs(left: int, right: int) -> int:
            if left >= right:
                return 0
            if right - left <= 2 and arr[left] == arr[right - 1]:
                return 1

            res = 0x7FFFFFFF
            # 枚举分割点
            for i in range(left, right):
                if arr[i] == arr[left]:

                    res = min(res, max(1, dfs(left + 1, i)) + dfs(i + 1, right))

            return res

        return dfs(0, len(arr))
```

2. **一个参数是端点，另一个参数是题目中的 target**型
   `关键：枚举分割点，left+1`
   `1335. 工作计划的最低难度.py`
   `1478. 安排邮筒.py`

```Python
 def dfs(cur: int, remain: int) -> int:
    if remain == 1:
        return ...

    res = 0x7FFFFFFF
    # 枚举分割点
    for i in range(cur + 1, len(houses)):
        res = min(res, calDistance(cur, i - 1) + dfs(i, remain - 1))
    return res
```
