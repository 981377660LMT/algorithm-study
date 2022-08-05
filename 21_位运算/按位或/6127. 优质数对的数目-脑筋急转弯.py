from bisect import bisect_left
from typing import List
from collections import defaultdict


MOD = int(1e9 + 7)
INF = int(1e20)


# !总结:
# !1. (a & b) << 1 + (a ^ b) = a + b
# !2. (a | b).bit_count() + (a & b).bit_count() = a.bit_count() + b.bit_count()
# !3. (a & b) + (a | b) = a + b
# !4. 前缀和/二分:固定某个数，求 sum(某个数>=k) 的个数

# 脑筋急转弯 + 两数之和


class Solution:
    def countExcellentPairs(self, nums: List[int], k: int) -> int:
        visited = set(nums)
        counter = defaultdict(int)
        for num in visited:
            counter[num.bit_count()] += 1

        # 前缀和查找
        preSum = [0]
        for i in range(70):
            preSum.append(preSum[-1] + counter[i])

        res = 0
        for cur in range(70):
            count1 = counter[cur]
            count2 = preSum[-1] - preSum[max(0, k - cur)]
            res += count1 * count2
        return res

    def countExcellentPairs2(self, nums: List[int], k: int) -> int:
        nums = list(set(nums))
        counts = []
        for num in nums:
            counts.append(num.bit_count())
        counts.sort()

        # 二分查找也可以
        res = 0
        for count in counts:
            pos = bisect_left(counts, k - count)
            res += len(nums) - pos
        return res


# print(Solution().countExcellentPairs(nums=[1, 2, 3, 1], k=3))
# print(Solution().countExcellentPairs(nums=[5, 1, 1], k=10))
# print(Solution().countExcellentPairs(nums=[1, 2, 3, 1, 536870911], k=3))
