from itertools import groupby, pairwise

# 统计并返回具有相同数量 0 和 1 的非空（连续）子字符串的数量，
# 并且这些子字符串中的所有 0 和所有 1 都是成组连续的。
class Solution:
    def countBinarySubstrings(self, s: str) -> int:
        groups = [[char, len(list(group))] for char, group in groupby(s)]
        return sum(min(pre, cur) for (_, pre), (_, cur) in pairwise(groups))

