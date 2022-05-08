from functools import lru_cache

MOD = int(1e9 + 7)


@lru_cache(None)
def cal(i: int, j: int) -> int:
    """第一类斯特林数:i个人,j个圆排列
    
    i,j<=1000
    - 将该新元素置于一个单独的轮换中 `cal(i - 1, j - 1)
    - 将该元素插入到任何一个现有的轮换中 `(i - 1) * cal(i - 1, j)``
    """
    if i == 0:
        return int(j == 0)
    return (cal(i - 1, j - 1) + (i - 1) * cal(i - 1, j)) % MOD


class Solution:
    def rearrangeSticks(self, n: int, k: int) -> int:
        """长度为从 1 到 n 的整数。请你将这些木棍排成一排，并满足从左侧 可以看到 恰好 k 根木棍
        
        划分为k个部分,每个部分排列种数为圆排列数(每个部分的最大值站在开头)
        """
        return cal(n, k)


print(Solution().rearrangeSticks(n=3, k=2))  # 5
