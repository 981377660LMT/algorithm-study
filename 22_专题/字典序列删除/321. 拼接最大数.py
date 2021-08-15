class Solution:
    def maxNumber(self, nums1, nums2, k):
        def pick_max(nums, k):
            stack = []
            drop = len(nums) - k
            for num in nums:
                while drop and stack and stack[-1] < num:
                    stack.pop()
                    drop -= 1
                stack.append(num)
            return stack[:k]

        def merge(A, B):
            ans = []
            while A or B:
                bigger = A if A > B else B
                ans.append(bigger[0])
                bigger.pop(0)
            return ans

        return max(
            merge(pick_max(nums1, i), pick_max(nums2, k - i))
            for i in range(k + 1)
            if i <= len(nums1) and k - i <= len(nums2)
        )

