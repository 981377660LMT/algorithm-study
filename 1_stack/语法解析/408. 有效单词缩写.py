# 字符串可以用 缩写 进行表示，缩写 的方法是将任意数量的 不相邻 的子字符串替换为相应子串的长度
# 0. 记录数字长度和字符下标
# 1. 数字没有'00'这种情况
# 2. 字符下标加上数字长度必须要对应 且不超出范围
# 3. 字符下标加数字长度最后等于word长度
class Solution:
    def validWordAbbreviation(self, word: str, abbr: str) -> bool:
        n = len(word)
        count = 0
        word_idx = 0
        for char in abbr:
            if char.isdigit():
                if count == 0 and char == '0':
                    return False
                count = count * 10 + int(char)
            else:
                word_idx += count
                if word_idx >= n or word[word_idx] != char:
                    return False
                count = 0
                word_idx += 1

        return word_idx + count == n


# 输入：word = "internationalization", abbr = "i12iz4n"
# 输出：true
# 解释：单词 "internationalization" 可以缩写为 "i12iz4n" ("i nternational iz atio n") 。


# 输入：word = "apple", abbr = "a2e"
# 输出：false
# 解释：单词 "apple" 无法缩写为 "a2e" 。

