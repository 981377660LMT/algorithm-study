from typing import List
from collections import Counter

# 1 <= nums.length <= 105

# 每一步操作中，你需要从数组中选出和为 k 的两个整数，并将它们移出数组。
# 返回你可以对数组执行的最大操作数。


class Solution:
    def maxOperations(self, nums: List[int], k: int) -> int:
        n = len(nums)
        nums.sort()
        res = 0
        L, R = 0, n - 1
        while L < R:
            if nums[L] + nums[R] == k:
                res += 1
                L += 1  # 2个数删除
                R -= 1
            elif nums[L] + nums[R] < k:
                L += 1  # 让小的大一些
            else:
                R -= 1  # 让大的小一些
        return res

    def maxOperations2(self, nums: List[int], k: int) -> int:
        freq = Counter(nums)
        res = 0
        for key in freq.keys():
            if key * 2 == k:
                res += freq[key] // 2
            # 只算一次，算小的
            elif key < k - key and k - key in freq:
                res += min(freq[key], freq[k - key])
        return res


print(Solution().maxOperations(nums=[1, 2, 3, 4], k=5))
# 输出：2
# 解释：开始时 nums = [1,2,3,4]：
# - 移出 1 和 4 ，之后 nums = [2,3]
# - 移出 2 和 3 ，之后 nums = []
# 不再有和为 5 的数对，因此最多执行 2 次操作。
