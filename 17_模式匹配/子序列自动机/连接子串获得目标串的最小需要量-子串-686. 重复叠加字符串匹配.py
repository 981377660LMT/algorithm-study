# 给定两个字符串 a 和 b，寻找重复叠加字符串 a 的最小次数，
# 使得字符串 b 成为叠加后的字符串 a 的子串，如果不存在则返回 -1。


# 注意：字符串 "abc" 重复叠加 0 次是 ""，重复叠加 1 次是 "abc"，
# 重复叠加 2 次是 "abcabc"。


# !覆盖b最少需要ceil(b/a)个a,最多需要ceil(b/a)+1个a
# !否则不能
class Solution:
    def repeatedStringMatch(self, a: str, b: str) -> int:
        na, nb = len(a), len(b)
        ceiling = (nb + na - 1) // na
        if b in a * ceiling:
            return ceiling
        if b in a * (ceiling + 1):
            return ceiling + 1
        return -1
