from typing import List
from itertools import accumulate
from collections import Counter

# 请你返回所有下标对 0 <= i, j < nums.length 的 floor(nums[i] / nums[j]) 结果之和。
# 1 <= nums.length <= 105

# 前缀和+枚举分母

# 1.枚举除数divisor
# 2.枚举divisor整数倍的被除数dividend，在这些点处对差分数组+1(有counter[divisor]组)
# 3.对差分数组求前缀和，可以得到每个整数作为被除数的贡献
# 4.对nums里的所有数，求出它们的贡献之和

MOD = int(1e9 + 7)


class Solution:
    def sumOfFlooredPairs(self, nums: List[int]) -> int:
        diff = [0] * (max(nums) + 1)
        counter = Counter(nums)

        # 枚举分母
        for fenmu in counter:
            # 枚举除数的贡献
            for fenzi in range(fenmu, len(diff), fenmu):
                diff[fenzi] += counter[fenmu]
        contribution = list(accumulate(diff))

        # 每个数作为分子的贡献
        return sum(contribution[num] for num in nums) % MOD


print(Solution().sumOfFlooredPairs(nums=[2, 5, 9]))
# 输出：10
# 解释：
# floor(2 / 5) = floor(2 / 9) = floor(5 / 9) = 0
# floor(2 / 2) = floor(5 / 5) = floor(9 / 9) = 1
# floor(5 / 2) = 2
# floor(9 / 2) = 4
# floor(9 / 5) = 1
# 我们计算每一个数对商向下取整的结果并求和得到 10 。

# 分數 的英文是 fraciton ，上面的 分子 是 numerator ，下面的 分母 是 denominator 。

