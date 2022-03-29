# return the number of k-length substrings that occur more than once in s.
from collections import defaultdict


class Solution:
    def solve(self, s, k):
        counter = defaultdict(int)

        for r in range(k, len(s) + 1):
            counter[s[r - k : r]] += 1

        return sum(1 for string in counter if counter[string] > 1)
