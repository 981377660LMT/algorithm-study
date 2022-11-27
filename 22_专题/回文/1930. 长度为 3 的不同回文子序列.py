from string import ascii_lowercase

# !3 <= s.length <= 105
# 你一个字符串 s ，返回 s 中 长度为 3 的不同回文子序列 的个数。
# Time O(26n)

# 寻找两端相同的位置


class Solution:
    def countPalindromicSubsequence(self, s: str) -> int:
        res = 0
        for char in ascii_lowercase:
            # 首次和末次出现的位置
            left, right = s.find(char), s.rfind(char)
            if left > -1:
                res += len(set(s[left + 1 : right]))
        return res


print(Solution().countPalindromicSubsequence(s="aabca"))
# 输出：3
# 解释：长度为 3 的 3 个回文子序列分别是：
# - "aba" ("aabca" 的子序列)
# - "aaa" ("aabca" 的子序列)
# - "aca" ("aabca" 的子序列)
