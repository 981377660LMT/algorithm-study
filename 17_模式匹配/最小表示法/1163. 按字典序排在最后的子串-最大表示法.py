# 1 <= s.length <= 4 * 10^5
class Solution:
    def lastSubstring2(self, s: str) -> str:
        res = ""
        for i in range(len(s)):
            res = max(res, s[i:])
        return res

    # '最大表示法' i,j 快慢指针，k表示i,j开头的字符串相等的长度
    def lastSubstring(self, s: str) -> str:
        n = len(s)
        i, j, k = 0, 1, 0
        while j + k < n:
            if s[i + k] == s[j + k]:
                k += 1
                continue
            # j没用了，要换掉
            elif s[i + k] > s[j + k]:
                j += k + 1
            else:
                i += k + 1
            if i == j:
                j += 1
            k = 0
        return s[i:]


print(Solution().lastSubstring("abab"))
# 输出："bab"
# 解释：我们可以找出 7 个子串 ["a", "ab", "aba", "abab", "b", "ba", "bab"]。按字典序排在最后的子串是 "bab"。
res = 0
