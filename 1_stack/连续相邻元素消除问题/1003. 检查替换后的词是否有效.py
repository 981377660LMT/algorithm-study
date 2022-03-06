# 碰到abc就消 最后能不能变0 => stack 消消乐
class Solution:
    def isValid(self, s: str) -> bool:
        stack = []
        for char in s:
            if char == 'c' and len(stack) >= 2 and stack[-1] == 'b' and stack[-2] == 'a':
                stack.pop()
                stack.pop()
            else:
                stack.append(char)

        return not stack


print(Solution().isValid(s="aabcbc"))
