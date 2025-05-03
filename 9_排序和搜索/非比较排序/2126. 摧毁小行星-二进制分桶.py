"""二进制分桶O(n)"""

# https://leetcode.cn/problems/destroying-asteroids/description/
# 如果行星碰撞时的质量 大于等于 小行星的质量，那么小行星被 摧毁 ，
# 并且行星会 获得 这颗小行星的质量。否则，行星将被摧毁。
#
# 按顺序考虑所有非空的组,如果当前mass小于组内最小值,
# 那么答案是false;如果当前mass大于等于组内最小值,
# !那么加上最小值之后必然大于组内所有值.
# !把所有数分到O(logU)个桶里，只检查桶中最小的元素.
#
# 1 <= asteroids.length <= 10^5
# 1 <= asteroids[i] <= 10^5

from typing import List


class Solution:
    def asteroidsDestroyed(self, mass: int, asteroids: List[int]) -> bool:
        max_ = max(asteroids)
        bitLen = max_.bit_length()
        bucketMins, bucketSums = [-1] * (bitLen + 1), [0] * (bitLen + 1)
        for num in asteroids:
            index = num.bit_length()
            bucketSums[index] += num
            bucketMins[index] = num if bucketMins[index] == -1 else min(num, bucketMins[index])

        for min_, sum_ in zip(bucketMins, bucketSums):
            if mass < min_:
                return False
            mass += sum_
        return True


print(Solution().asteroidsDestroyed(mass=10, asteroids=[3, 9, 19, 5, 21]))
