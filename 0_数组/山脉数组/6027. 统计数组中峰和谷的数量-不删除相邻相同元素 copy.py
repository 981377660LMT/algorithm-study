from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def countHillValley(self, nums: List[int]) -> int:
        stack = []
        for num in nums:
            if not stack or num != stack[-1][0]:
                stack.append((num, 1))
            else:
                stack[-1] = (num, stack[-1][1] + 1)

        res = 0
        for i in range(1, len(stack) - 1):
            if stack[i - 1][0] < stack[i][0] > stack[i + 1][0]:
                res += stack[i][1]
            elif stack[i - 1][0] > stack[i][0] < stack[i + 1][0]:
                res += stack[i][1]

        return res


print(Solution().countHillValley(nums=[2, 4, 1, 1, 6, 5]))
