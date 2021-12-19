# 只包含元音
class Solution:
    def longestBeautifulSubstring(self, word: str) -> int:
        res, length, type = 0, 1, 1
        for right in range(1, len(word)):
            if word[right] >= word[right - 1]:
                length += 1
            if word[right] > word[right - 1]:
                type += 1
            if word[right] < word[right - 1]:
                length, type = 1, 1
            if type == 5:
                res = max(res, length)

        return res


print(Solution().longestBeautifulSubstring(word="aeiaaioaaaaeiiiiouuuooaauuaeiu"))
# 输出：13
# 解释：最长子字符串是 "aaaaeiiiiouuu" ，长度为 13 。
