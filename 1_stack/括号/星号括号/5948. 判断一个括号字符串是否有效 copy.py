from typing import List


class Solution:
    def canBeValid(self, s: str, locked: str) -> bool:
        if len(s) & 1:
            return False
        n, sb = len(s), list(s)
        for i in range(n):
            if locked[i] == '0':
                sb[i] = '*'

        return self.check(sb, isReversed=False) and self.check(sb, isReversed=True)

    def check(self, brackets: List[str], isReversed: bool) -> bool:
        level = 0
        score = 1
        if isReversed:
            brackets = list(reversed(brackets))
            score = -1

        for char in brackets:
            if char == '(':
                level += score
            elif char == ')':
                level -= score
            elif char == '*':
                level += 1

            if level < 0:
                return False

        return True


print(Solution().canBeValid(s="))()))", locked="010100"))

