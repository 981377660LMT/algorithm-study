# n ≤ 100,000
# 找到最短子串,使得其中一个字母freq大于子串内其他字母freq之和
class Solution:
    def solve(self, s):
        n = len(s)
        if any(s[i] == s[i + 1] for i in range(n - 1)):
            return 2
        if any(s[i] == s[i + 2] for i in range(n - 2)):
            return 3
        return -1

