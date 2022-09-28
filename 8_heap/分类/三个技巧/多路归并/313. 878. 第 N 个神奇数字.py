# 如果正整数可以被 A 或 B 整除，那么它是神奇的。
# 返回第 N 个神奇数字。由于答案可能非常大，返回它模 10^9 + 7 的结果。
from math import lcm

MOD = int(1e9 + 7)


class Solution:
    def nthMagicalNumber(self, n: int, a: int, b: int) -> int:
        def check(mid: int) -> bool:
            """[1,mid]中神奇数>=n"""
            total = mid // a + mid // b - mid // lcm_
            return total >= n

        lcm_ = lcm(a, b)
        left, right = 1, int(1e18)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left % MOD


print(Solution().nthMagicalNumber(n=4, a=2, b=3))
# 输出：6
