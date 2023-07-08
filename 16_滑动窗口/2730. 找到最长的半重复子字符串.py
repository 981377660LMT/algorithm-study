# 如果一个字符串 t 中至多有一对相邻字符是相等的，那么称这个字符串 t 是 半重复的 。
# 例如，0010 、002020 、0123 、2002 和 54944 是半重复字符串，
# 而 00101022 和 1101234883 不是。


# 2730. 找到最长的半重复子字符串
# https://leetcode.cn/problems/find-the-longest-semi-repetitive-substring/
class Solution:
    def longestSemiRepetitiveSubstring(self, s: str) -> int:
        res, left, n = 0, 0, len(s)
        same = 0
        for right in range(n):
            same += right >= 1 and s[right] == s[right - 1]
            while left <= right and same > 1:
                same -= left + 1 < n and s[left] == s[left + 1]
                left += 1
            res = max(res, right - left + 1)
        return res


if __name__ == "__main__":
    assert Solution().longestSemiRepetitiveSubstring("0001") == 3
