from typing import List
from SA import useSA
from atcSA import sa_is, rank_lcp


class Solution:
    def countDistinct(self, nums: List[int], k: int, p: int) -> int:
        # https://leetcode.cn/problems/k-divisible-elements-subarrays/solution/by-freeyourmind-2m6j/
        right = count = res = 0
        rightMosts = [0] * (n := len(nums))
        for left, num in enumerate(nums):
            # 先用双指针O(n)的时间计算出所有满足条件的子数组的数量 注意要枚举后缀(固定left 移动right)
            while right < n and count + (mod_ := nums[right] % p == 0) <= k:
                count += mod_
                right += 1
            res += right - left
            rightMosts[left] = right
            count -= num % p == 0
        sa = sa_is(nums, max(nums))
        _, lcp = rank_lcp(nums, sa)
        lcp = [0] + lcp
        return res - sum(min(lcp[i], rightMosts[sa[i]] - sa[i]) for i in range(n))


print(Solution().countDistinct(nums=[2, 3, 3, 2, 2], k=2, p=2))
print(Solution().countDistinct(nums=[6, 20, 5, 18], k=3, p=14))
