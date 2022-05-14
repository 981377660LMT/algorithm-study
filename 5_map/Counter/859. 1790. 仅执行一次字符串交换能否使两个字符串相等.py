# 如果对 其中一个字符串 执行 最多一次字符串交换 就可以使两个字符串相等，返回 true ；否则，返回 false 。
# s1.length == s2.length
class Solution:
    def areAlmostEqual(self, s1: str, s2: str) -> bool:
        """仅执行一次交换能否使两个字符串相等"""
        if s1 == s2:
            return True

        # 挑出不同的字符对,对数只能为2，并且对称，如 (a,b)与(b,a)
        diff = [(a, b) for a, b in zip(s1, s2) if a != b]
        return len(diff) == 2 and diff[0] == diff[1][::-1]
