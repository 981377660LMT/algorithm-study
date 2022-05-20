from collections import deque
from typing import List


class Solution:
    def maxNumber(self, nums1, nums2, k):
        def pickMax(nums: List[int], k: int) -> List[int]:
            """选k个数拼成最大数  栈顶肯定要最大 单减的单调栈"""
            stack = []
            drop = len(nums) - k
            for num in nums:
                while drop and stack and stack[-1] < num:
                    stack.pop()
                    drop -= 1
                stack.append(num)
            return stack[:k]

        def merge(A: List[int], B: List[int]) -> List[int]:
            """合并两个数组，字典序最大"""
            res = []
            nums1, nums2 = deque(A), deque(B)
            while nums1 or nums2:
                bigger = nums1 if nums1 > nums2 else nums2
                res.append(bigger.popleft())
            return res

        return max(
            merge(pickMax(nums1, i), pickMax(nums2, k - i))
            for i in range(k + 1)
            if i <= len(nums1) and k - i <= len(nums2)
        )

