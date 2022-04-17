import string


allChar = string.ascii_uppercase
digitByChar = {char: i + 1 for i, char in enumerate(allChar)}
RADIX = 26


class Solution:
    def titleToNumber(self, columnTitle: str) -> int:
        res = 0
        base = 1
        for i in range(len(columnTitle) - 1, -1, -1):
            char = columnTitle[i]
            res += base * digitByChar[char]
            base *= RADIX
        return res
