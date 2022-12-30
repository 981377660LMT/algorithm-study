# 2518. 好分区的数目
# 给你一个正整数数组 nums 和一个整数 k 。
# 分区 的定义是：将数组划分成两个有序的 组 ，并满足每个元素 恰好 存在于 某一个 组中。
# 如果分区中每个组的元素和都大于等于 k ，则认为分区是一个好分区。
# 返回 不同 的好分区的数目。由于答案可能很大，请返回对 109 + 7 取余 后的结果。
# 如果在两个分区中，存在某个元素 nums[i] 被分在不同的组中，则认为这两个分区不同。
# 1 <= nums.length, k <= 1000
# 1 <= nums[i] <= 1e9


from typing import List
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def countPartitions(self, nums: List[int], k: int) -> int:
        def cal(upper: int) -> int:
            """选取若干个数使得和严格小于upper"""
            dp = defaultdict(int, {0: 1})
            for num in nums:
                ndp = dp.copy()
                for key, val in dp.items():
                    if key + num < upper:
                        ndp[key + num] = (ndp[key + num] + val) % MOD
                dp = ndp

            res = 0
            for key, val in dp.items():
                res = (res + val) % MOD
            return res

        if sum(nums) < 2 * k:
            return 0
        n = len(nums)
        return (pow(2, n, MOD) - 2 * cal(k)) % MOD


print(Solution().countPartitions(nums=[1, 2, 3, 4], k=4))
print(Solution().countPartitions(nums=[3, 3, 3], k=4))
print(Solution().countPartitions(nums=[6, 6], k=2))
print(Solution().countPartitions(nums=[1, 1, 1, 1], k=3))
