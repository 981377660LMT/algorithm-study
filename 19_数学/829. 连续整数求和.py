# 给定一个正整数 N，试求有多少组连续`正整数`满足所有数字之和为 N?
# !等差数列求和/连续(正)整数求和


class Solution:
    def consecutiveNumbersSum(self, n: int) -> int:
        """
        有多少组连续`正整数`满足所有数字之和为 N?
        a: 首项 c: 项数
        a*c + c*(c-1)/2 = n
        a*c = n - c*(c-1)/2 是整数
        """
        res = 0
        count = 1  # 枚举项数
        while count * (count - 1) // 2 < n:
            if (n - count * (count - 1) // 2) % count == 0:
                res += 1
            count += 1
        return res

    # !如果不要求正整数呢?
    def consecutiveNumbersSum2(self, n: int) -> int:
        """
        有多少组连续`整数`满足所有数字之和为 N
        乘以2是因为非正数也可以取到
        sum(left,right) = sum(-left+1,right)
        """
        return 2 * self.consecutiveNumbersSum(n)


# (2*a+count-1)*count=2*n
# a=n/c - (c-1)/2
print(Solution().consecutiveNumbersSum(n=9))
