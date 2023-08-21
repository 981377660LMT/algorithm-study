# 2831. 找出最长等值子数组
# https://leetcode.cn/problems/find-the-longest-equal-subarray/description/
# 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。
# 如果子数组中所有元素都相等，则认为子数组是一个 等值子数组 。注意，空数组是 等值子数组 。
# 从 nums 中删除最多 k 个元素后，返回可能的最长等值子数组的长度。
# 子数组 是数组中一个连续且可能为空的元素序列


# 1. 二分+定长滑动窗口
# 2. 哈希表+不定长滑动窗口


from collections import defaultdict
from typing import List


class Solution:
    def longestEqualSubarray(self, nums: List[int], k: int) -> int:
        """哈希表+不定长滑动窗口"""

        def cal(indexes: List[int]) -> int:
            res, left = 0, 0
            for right in range(len(indexes)):
                while left <= right and indexes[right] - indexes[left] + 1 > k + right - left + 1:
                    left += 1
                res = max(res, right - left + 1)
            return res

        mp = defaultdict(list)
        for i, v in enumerate(nums):
            mp[v].append(i)

        res = 0
        for v in mp.values():
            res = max(res, cal(v))
        return res

    def longestEqualSubarray2(self, nums: List[int], k: int) -> int:
        """二分+定长滑动窗口"""

        def check(mid: int) -> bool:
            """是否存在长度不超过mid+k的子数组,某个数的频率>=mid"""
            n = len(nums)
            counter = defaultdict(int)
            winLen = mid + k
            for right in range(n):
                counter[nums[right]] += 1
                if right >= winLen:
                    counter[nums[right - winLen]] -= 1
                if counter[nums[right]] >= mid:
                    return True
            return False

        n = len(nums)
        left, right = 0, n
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


if __name__ == "__main__":
    print(Solution().longestEqualSubarray([1], 1))
    print(Solution().longestEqualSubarray2([1, 2, 1], 0))
