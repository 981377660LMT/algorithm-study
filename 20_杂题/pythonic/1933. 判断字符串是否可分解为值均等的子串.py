from itertools import groupby


# 一个字符串的所有字符都是一样的，被称作等值字符串。
# 如果有且仅有一个等值子字符串长度为2，其他的等值子字符串的长度都是3.就返回真
class Solution:
    def isDecomposable(self, s: str) -> bool:
        match = set(len(list(g)) % 3 for _, g in groupby(s))
        return sum(match) == 2 and 1 not in match


print([list(g) for e, g in groupby("00011111222")])


