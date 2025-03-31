# 勒让德定理
# P在N阶乘中出现的次数就是看1-N中所有数字
# 出现P因子的次数总和
# 根据勒让德定理就是 N//P + N//(P^2) + N//(p^3) + ......
# 既可以每次除以5 然后累加每次除后的 n

# 阶乘后的0


# 0 <= n <= 104
class Solution:
    def trailingZeroes(self, n: int) -> int:
        """勒让德定理求n!的尾随0的个数,时间复杂度O(log(n))"""
        res = 0
        while n:
            n //= 5
            res += n
        return res

    def trailingZeroes2(self, n: int) -> int:
        """枚举因子求n!的尾随0的个数,时间复杂度O(n)"""
        res = 0
        for i in range(5, n + 1, 5):
            while i % 5 == 0:
                i //= 5
                res += 1
        return res
