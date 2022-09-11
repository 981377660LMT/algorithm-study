# 返回 重新排列空格后的字符串 。
# 请你重新排列空格，使每对相邻单词之间的空格数目都 相等 ，并尽可能 最大化 该数目。
# 如果不能重新平均分配所有空格，请 将多余的空格放置在字符串末尾 ，

# !divmod平均分组


class Solution:
    def reorderSpaces(self, text: str) -> str:
        white = text.count(" ")
        words = text.split()  # !默认分隔符为所有的空字符
        if len(words) == 1:
            return words[0] + " " * white
        div, mod = divmod(white, len(words) - 1)
        return (" " * div).join(words) + " " * mod


print(Solution().reorderSpaces("  this   is  a sentence "))

# 输出："this   is   a   sentence"
# 解释：总共有 9 个空格和 4 个单词。可以将 9 个空格平均分配到相邻单词之间，相邻单词间空格数为：9 / (4-1) = 3 个。
