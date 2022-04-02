class Solution:
    def solve(self, nums):
        stack = []
        res = 0
        for num in nums:
            while len(stack) > 0 and stack[-1] > num:
                stack.pop()
            stack.append(num)
            res += len(stack)
        return res

