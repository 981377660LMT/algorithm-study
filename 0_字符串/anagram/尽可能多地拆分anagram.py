from collections import Counter

# 尽可能多地拆分anagram
# Anagram Partitioning


class Solution:
    def solve(self, a, b):
        """create a cut as soon as we find that the substrings we are considering are anagrams"""
        cut = []
        start = 0
        counter = Counter()
        for i, (x, y) in enumerate(zip(a, b)):
            counter[x] += 1
            if counter[x] == 0:
                del counter[x]
            counter[y] -= 1
            if counter[y] == 0:
                del counter[y]

            if not counter:
                cut.append(start)
                start = i + 1
        return cut if start == len(a) else []


print(Solution().solve(a="catdogwolf", b="actgodflow"))
# [0, 2, 3, 6]
