# 第二类斯特林数
from functools import lru_cache

MOD = int(1e9 + 7)


@lru_cache(None)
def cal(i: int, j: int) -> int:
    """第二类斯特林数:i个人,j个子集
    
    i,j<=1000
    - 将新元素单独放入一个子集 `cal(i - 1, j - 1)`
    - 将新元素放入一个现有的非空子集 `j * cal(i - 1, j)`
    """
    if i == 0:
        return int(j == 0)
    return (cal(i - 1, j - 1) + j * cal(i - 1, j)) % MOD


# dp[i][j]表示i个盒子 j颗糖
class Solution:
    def waysToDistribute2(self, n: int, k: int) -> int:
        dp = [[0] * (n + 1) for _ in range(k + 1)]
        for i in range(1, k + 1):
            dp[i][i] = 1
        for i in range(1, k + 1):
            for j in range(i + 1, n + 1):
                # 新的糖独占1盒 dp[i-1][j-1]
                # 不独占一盒随意放 i*dp[i][j - 1]
                dp[i][j] = (dp[i - 1][j - 1] + i * dp[i][j - 1]) % int(1e9 + 7)
        return dp[-1][-1]

    def waysToDistribute(self, n: int, k: int) -> int:
        return cal(n, k)


print(Solution().waysToDistribute(n=4, k=2))
# 输出：7
# 解释：把糖果 4 分配到 2 个手袋中的一个，共有 7 种方式:
# (1), (2,3,4)s
# (1,2), (3,4)
# (1,3), (2,4)
# (1,4), (2,3)
# (1,2,3), (4)
# (1,2,4), (3)
# (1,3,4), (2)
