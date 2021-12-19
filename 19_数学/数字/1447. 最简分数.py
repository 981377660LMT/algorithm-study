# 1 <= n <= 100
from typing import List
from math import gcd


# 请你返回所有 0 到 1 之间（不包括 0 和 1）满足分母小于等于  n 的 最简 分数
class Solution:
    def simplifiedFractions(self, n: int) -> List[str]:
        return [f"{j}/{i}" for i in range(2, n + 1) for j in range(1, i) if gcd(i, j) == 1]


print(Solution().simplifiedFractions(4))
# 输入：n = 4
# 输出：["1/2","1/3","1/4","2/3","3/4"]
# 解释："2/4" 不是最简分数，因为它可以化简为 "1/2" 。
