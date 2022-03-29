# return whether it is possible to cut both strings at a common point such that the first part of a and the second part of b form a palindrome
# 字符串a与字符串b头部交换，问是否可以形成回文串

# 回文只需要检查i<n//2 的部分
class Solution:
    def solve(self, a: str, b: str) -> bool:
        def check(s1: str, s2: str) -> bool:
            i = 0
            while i < n // 2 and s1[i] == s2[~i]:
                i += 1
            return all(s1[j] == s1[~j] for j in range(i, n // 2))

        n = len(a)
        return check(a, b) or check(b, a)

