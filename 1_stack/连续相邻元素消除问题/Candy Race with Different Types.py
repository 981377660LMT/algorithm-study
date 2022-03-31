# 不断选择相邻元素消除，不能继续就输了
class Solution:
    def solve(self, nums):
        """return whether you will win."""
        stack = []
        count = 0
        for num in nums:
            if stack and stack[-1] == num:
                stack.pop()
                count += 1
            else:
                stack.append(num)

        return not not count & 1

