from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minDeletion(self, nums: List[int]) -> int:
        stack = []
        for num in nums:
            if not stack:
                stack.append(num)
                continue
            if len(stack) & 1:
                if stack[-1] != num:
                    stack.append(num)
                else:
                    continue
            else:
                stack.append(num)

        res = len(stack)
        cand = len(nums) - res
        return cand + 1 if res & 1 else cand


print(Solution().minDeletion(nums=[1, 1, 2, 3, 5]))
print(Solution().minDeletion([1, 1, 2, 2, 3, 3]))

