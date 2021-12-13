# 请你将其中 一个 字符用任意小写英文字母替换，使得结果字符串的 字典序最小 ，且 不是 回文串。

# Check half of the string,
# replace a non 'a' character to 'a'.

# If only one character, return empty string.
# Otherwise repalce the last character to 'b'
class Solution:
    def breakPalindrome(self, palindrome: str) -> str:
        if len(palindrome) <= 1:
            return ''
        half = len(palindrome) >> 1
        for i in range(half):
            if palindrome[i] != 'a':
                return palindrome[:i] + 'a' + palindrome[i + 1 :]
        return palindrome[:-1] + 'b'


print(Solution().breakPalindrome(palindrome="abccba"))
# 输出："aaccba"
# 解释：存在多种方法可以使 "abccba" 不是回文，例如 "zbccba", "aaccba", 和 "abacba" 。
# 在所有方法中，"aaccba" 的字典序最小。

