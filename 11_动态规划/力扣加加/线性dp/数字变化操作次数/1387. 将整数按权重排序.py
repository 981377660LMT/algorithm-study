from functools import lru_cache

# 我们将整数 x 的 权重 定义为按照下述规则将 x 变成 1 所需要的步数：
# 如果 x 是偶数，那么 x = x / 2
# 如果 x 是奇数，那么 x = 3 * x + 1

# 请你返回区间 [lo, hi] 之间的整数按权重排序后的第 k 个数。

# 注意要写在Solution外 避免实例化每次都重复计算
@lru_cache(None)
def dfs(x: int):
    if x == 1:
        return 0
    if not x & 1:
        return 1 + dfs(x // 2)
    return 1 + dfs(3 * x + 1)


class Solution:
    def getKth(self, lo: int, hi: int, k: int) -> int:
        return sorted(range(lo, hi + 1), key=dfs)[k - 1]


print(Solution().getKth(lo=12, hi=15, k=2))
# 输出：13
# 解释：12 的权重为 9（12 --> 6 --> 3 --> 10 --> 5 --> 16 --> 8 --> 4 --> 2 --> 1）
# 13 的权重为 9
# 14 的权重为 17
# 15 的权重为 17
# 区间内的数按权重排序以后的结果为 [12,13,14,15] 。对于 k = 2 ，答案是第二个整数也就是 13 。
# 注意，12 和 13 有相同的权重，所以我们按照它们本身升序排序。14 和 15 同理。
