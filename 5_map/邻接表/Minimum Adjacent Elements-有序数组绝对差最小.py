# we say that two numbers nums[i] ≤ nums[j] are adjacent if there's no number in between (nums[i], nums[j]) in nums.
# Return the minimum possible abs(j - i) such that nums[j] and nums[i] are adjacent.


# 注意到排序之后,adjacent元素相邻

from collections import defaultdict


class Solution:
    def solve(self, nums):
        indexMap = defaultdict(list)
        for i, num in enumerate(nums):
            indexMap[num].append(i)

        res = int(1e20)

        # 相等元素
        for indexes in indexMap.values():
            for pre, cur in zip(indexes, indexes[1:]):
                res = min(res, abs(pre - cur))

        if res == 1:
            return 1

        # 不等元素
        keys = sorted(indexMap)
        for i in range(len(keys) - 1):
            nums1, nums2 = indexMap[keys[i]], indexMap[keys[i + 1]]
            i, j = 0, 0
            while i < len(nums1) and j < len(nums2):
                res = min(res, abs(nums1[i] - nums2[j]))
                if nums1[i] < nums2[j]:
                    i += 1
                else:
                    j += 1

        return res


print(Solution().solve(nums = [0, -10, 5, -5, 1]))
