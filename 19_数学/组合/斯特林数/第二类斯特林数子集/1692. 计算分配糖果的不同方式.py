# 第二类斯特林数
from functools import lru_cache

MOD = int(1e9 + 7)


@lru_cache(None)
def cal2(i: int, j: int) -> int:
    """第二类斯特林数:i个人,j个子集
    
    i,j<=1000
    - 将新元素单独放入一个子集 `cal(i - 1, j - 1)`
    - 将新元素放入一个现有的非空子集 `j * cal(i - 1, j)`
    """
    if i == 0:
        return int(j == 0)
    return (cal2(i - 1, j - 1) + j * cal2(i - 1, j)) % MOD


dp2 = [[0] * 1001 for _ in range(1001)]
dp2[0][0] = 1
for i in range(1, 1001):
    for j in range(1, 1001):
        dp2[i][j] = (dp2[i - 1][j - 1] + j * dp2[i - 1][j]) % MOD

# dp[i][j]表示i个盒子 j颗糖
class Solution:
    def waysToDistribute(self, n: int, k: int) -> int:
        return dp2[n][k]
        return cal2(n, k)


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
