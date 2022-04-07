from collections import Counter


class Solution:
    def solve(self, a, b):
        return len(list(((Counter(a) & Counter(b)).elements())))
