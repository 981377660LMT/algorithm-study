# 给定一个正整数 N，试求有多少组连续正整数满足所有数字之和为 N?
# !等差数列求和/连续整数求和
EPS = 1e-6


class Solution:
    def consecutiveNumbersSum(self, n: int) -> int:
        """
        a: 首项 c: 项数
        (2*a+c-1)*c=2*n
        a=n/c - (c-1)/2
        """
        res = 0

        count = 1  # 枚举项数
        while count * (count - 1) // 2 < n:
            first = n / count - (count - 1) / 2  # 首项
            res += abs(first - int(first)) < EPS
            count += 1
        return res


# (2*a+count-1)*count=2*n
# a=n/c - (c-1)/2
print(Solution().consecutiveNumbersSum(n=9))
