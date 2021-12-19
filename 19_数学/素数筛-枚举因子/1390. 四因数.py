# 给你一个整数数组 nums，请你返回该数组中恰有四个因数的这些整数的各因数之和。
from typing import List
from math import floor, sqrt


class Solution:
    # O(n*根号c)
    def sumFourDivisors(self, nums: List[int]) -> int:
        res = 0
        for num in nums:
            divisor = set()
            for i in range(1, floor(sqrt(num)) + 1):
                if num % i == 0:
                    divisor.add(num // i)
                    divisor.add(i)
                if len(divisor) > 4:
                    break

            if len(divisor) == 4:
                res += sum(divisor)
        return res


print(Solution().sumFourDivisors(nums=[21, 4, 7]))
# 输出：32
# 解释：
# 21 有 4 个因数：1, 3, 7, 21
# 4 有 3 个因数：1, 2, 4
# 7 有 2 个因数：1, 7
# 答案仅为 21 的所有因数的和。
