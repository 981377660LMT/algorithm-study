# 给你一个只包含 '(' 和 ')' 的字符串，找出最长有效（格式正确且连续）括号子串的长度。


from itertools import groupby


class Solution:
    def longestValidParentheses(self, s: str) -> int:
        n = len(s)
        isBad = [0] * n
        stack = []
        for i in range(n):
            if s[i] == "(":
                stack.append(i)
            else:
                if stack:
                    stack.pop()
                else:
                    isBad[i] = 1
        for i in stack:
            isBad[i] = 1

        groups = [(char, len(list(group))) for char, group in groupby(isBad)]
        return max((b for a, b in groups if a == 0), default=0)


assert Solution().longestValidParentheses("(()") == 2
