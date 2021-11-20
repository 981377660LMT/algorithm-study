from itertools import zip_longest
from re import findall

# 请你将该字符串重新格式化，使得任意两个相邻字符的类型都不同。也就是说，字母后面应该跟着数字，而数字后面应该跟着字母。
# 如果无法按要求重新格式化，则返回一个 空字符串 。
class Solution:
    def reformat(self, s: str) -> str:
        alphas = findall(r'[a-z]', s)
        digits = findall(r'\d', s)
        if abs(len(alphas) - len(digits)) > 1:
            return ''
        short, long = sorted([alphas, digits], key=len)
        return ''.join(a + b for a, b in zip_longest(long, short, fillvalue=''))


print(Solution().reformat("covid2019"))
