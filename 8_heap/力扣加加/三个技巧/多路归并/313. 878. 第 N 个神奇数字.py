# 如果正整数可以被 A 或 B 整除，那么它是神奇的。
from math import gcd

# 返回第 N 个神奇数字。由于答案可能非常大，返回它模 10^9 + 7 的结果。
class Solution:
    def nthMagicalNumber(self, n: int, a: int, b: int) -> int:
        def is_enough(mid: int) -> bool:
            total = mid // a + mid // b - mid // ab
            return total >= n

        # 计算最小公倍数
        ab = a * b // gcd(a, b)

        left, right = 1, 10 ** 16
        while left <= right:
            mid = (left + right) >> 1
            if is_enough(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left % (10 ** 9 + 7)


print(Solution().nthMagicalNumber(n=4, a=2, b=3))
# 输出：6
