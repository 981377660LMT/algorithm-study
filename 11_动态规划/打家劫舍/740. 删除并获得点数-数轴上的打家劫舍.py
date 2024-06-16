from typing import List
from collections import defaultdict


def max2(a: int, b: int) -> int:
    return a if a > b else b


def rob2(nums: List[int], gap=1) -> int:
    """给定n个数.选择数x后,不能选择[x-gap,x+gap]之间的数,求最大和."""
    if not nums:
        return 0
    counter = defaultdict(int)
    for v in nums:
        counter[v] += 1
    keys = sorted(counter)
    dp = [0] * (len(keys) + 1)
    left = 0
    for i, v in enumerate(keys):
        while keys[left] < v - gap:
            left += 1
        dp[i + 1] = max2(dp[i], dp[left] + v * counter[v])
    return dp[-1]


if __name__ == "__main__":

    class Solution:
        # 100316. 施咒的最大总伤害
        # https://leetcode.cn/problems/maximum-total-damage-with-spell-casting/description/
        def maximumTotalDamage(self, power: List[int]) -> int:
            return rob2(power, 2)

        # 740. 删除并获得点数
        # https://leetcode.cn/problems/delete-and-earn/
        def deleteAndEarn(self, nums: List[int]) -> int:
            return rob2(nums)
