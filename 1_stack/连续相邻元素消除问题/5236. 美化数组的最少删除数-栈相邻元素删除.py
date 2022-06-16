from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minDeletion(self, nums: List[int]) -> int:
        stack = []
        for num in nums:
            if stack and stack[-1] == num and ((len(stack) & 1)):
                continue
            stack.append(num)

        len_ = len(stack)
        cand = len(nums) - len_
        return cand + 1 if len_ & 1 else cand


print(Solution().minDeletion(nums=[1, 1, 2, 3, 5]))
print(Solution().minDeletion([1, 1, 2, 2, 3, 3]))

