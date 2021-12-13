# 给定一个正整数 a，找出最小的正整数 b 使得 b 的所有数位相乘恰好等于 a。
# 如果不存在这样的结果或者结果不是 32 位有符号整数，返回 0。

# 1.贪心---从9开始能分解就分解 最小到1
# 2.大的数字靠右，小的数字靠左 因为题目要求数字尽量小
# https://leetcode-cn.com/problems/minimum-factorization/solution/c-python3-tan-xin-neng-fen-jie-jiu-fen-j-1zzq/
class Solution:
    def smallestFactorization(self, num: int) -> int:
        if num <= 9:
            return num

        res = 0
        base = 1
        for factor in range(9, 1, -1):
            while num % factor == 0:
                num //= factor
                res = factor * base + res
                base *= 10

                if res > 2 ** 31 - 1:
                    return 0

        if num != 1:
            return 0

        return res


print(Solution().smallestFactorization(48))
# 68
