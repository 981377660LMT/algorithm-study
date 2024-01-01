# Leetcode题解：1574. 删除最短的子数组使剩余数组有序
# https://leetcode.cn/problems/shortest-subarray-to-be-removed-to-make-array-sorted/
# 删除子数组
# 请你删除一个子数组（可以为空），使得 arr 中剩下的元素是 非递减 的。
# 返回满足题目要求删除的最短子数组的长度。
# 同 https://leetcode.cn/problems/count-the-number-of-incremovable-subarrays-ii/

# 答案肯定是前缀+后缀组成的
# 处理出包含每个位置的前缀和后缀的非递减长度
# 然后滑窗找到每个左端点对应的右端点

from typing import List


class Solution:
    def findLengthOfShortestSubarray(self, arr: List[int]) -> int:
        # 找到连续递增的前缀和后缀的位置
        n = len(arr)
        i, j = 0, len(arr) - 1
        while i + 1 < n and arr[i] <= arr[i + 1]:
            i += 1
        if i == n - 1:  # 全部递增
            return 0
        while j - 1 >= 0 and arr[j] >= arr[j - 1]:
            j -= 1

        # !注意保留的前缀长为0和后缀长为0的情况要特判
        # !(这里必须特判,因为没有对应的arr[i])
        res = min(n - i - 1, j)

        # !1. 双指针
        # right = j
        # for left in range(i + 1):
        #     while right < n and arr[left] > arr[right]:
        #         right += 1
        #     res = min(res, right - left - 1)

        # !2. 二分
        for v in range(i + 1):
            # !找到最左边的left,使得nums[left] > nums[v]
            left, right = j, n - 1
            ok = False
            while left <= right:
                mid = (left + right) // 2
                if arr[mid] >= arr[v]:
                    right = mid - 1
                    ok = True
                else:
                    left = mid + 1
            if ok:
                res = min(res, left - v - 1)

        return res


print(Solution().findLengthOfShortestSubarray([1, 2, 3, 10, 4, 2, 3, 5]))
# [5,4,3,2,1]
# print(Solution().findLengthOfShortestSubarray([5, 4, 3, 2, 1]))
# print(Solution().findLengthOfShortestSubarray([1, 2, 3]))
# [16,10,0,3,22,1,14,7,1,12,15]
