# 20. 有效的括号
# https://leetcode.cn/problems/valid-parentheses/

MP = {"(": ")", "[": "]", "{": "}"}


class Solution:
    def isValid(self, s: str) -> bool:
        if len(s) & 1:
            return False
        stack = []
        for c in s:
            if c in MP:
                stack.append(c)
            elif not stack or MP[stack.pop()] != c:
                return False
        return not stack
