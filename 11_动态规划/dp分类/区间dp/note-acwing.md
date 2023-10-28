1. 环形石子合并
   `同一个石子 只会被合并到 一个石子堆` 里
   区间 DP 我们是把区间 [a,b]拆分为 [a,k]和[k+1,b]
2. 能量项链
   合并魔法石时，分割点 k 要被分到 左侧石子堆的右端点 和 右侧石子堆的左端点 中
   因此，`参数 k 要作为两个区间的共同端点来使用`，即 [a,k] 和 [k,b]
3. 加分二叉树
   后序 dfs+讨论根节点位置
4. 凸多边形的划分
5. 棋盘分割
   `n<15【记忆化搜索/二维区间 DP】`

环形问题：
`翻倍+从0到n求最值`

关键：
`left+right+枚举分割点`

**注意**
如果左右两端区间的代价不能共用中间点 k
就是 `dfs(left, k) + dfs(k + 1, right) + cost(left,k,right)`
否则为 `dfs(left, k) + dfs(k, right) + cost(left,k,right)`

```Python
@lru_cache(maxsize=None)
def dfs(left: int, right: int) -> int:
    """[left,right]这一段合并的代价最小"""
    # if right-left <= ...
    #     return 0
    if left >= right:
        return 0
    res = INF
    for i in range(left, right):
        res = min(res, dfs(left, i) + dfs(i + 1, right) + preSum[right + 1] - preSum[left])
    return res


print(dfs(0, n - 1))
```
