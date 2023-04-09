# 如果 n `恰好`有三个正除数 ，返回 true
# i*i =n ，并且i是质数的时候，才是刚刚好三个除数
from math import sqrt


class Solution:
    def isThree(self, n: int) -> bool:
        def isprime(n: float):
            if n < 2:
                return False
            for factor in range(2, int(sqrt(n)) + 1):
                if root % factor == 0:
                    return False
            return True

        root = sqrt(n)
        if not root.is_integer():
            return False

        return isprime(root)


# 输入：n = 2
# 输出：false
# 解释：2 只有两个除数：1 和 2 。
print(Solution().isThree(81))
print(Solution().isThree(4))
