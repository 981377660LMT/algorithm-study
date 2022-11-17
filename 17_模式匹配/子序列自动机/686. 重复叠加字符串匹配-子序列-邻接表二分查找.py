from bisect import bisect_right
from collections import defaultdict

# 邻接表+二分判断子序列


class Solution:
    def solve(self, s, t):
        def match(hit: int):
            """判断子序列"""
            prePos = -1
            while hit < len(t):
                nextPos = bisect_right(indexMap[t[hit]], prePos)
                if nextPos == len(indexMap[t[hit]]):
                    break

                prePos = indexMap[t[hit]][nextPos]
                hit += 1

            return hit

        if not set(t) <= (set(s)):
            return -1

        indexMap = defaultdict(list)
        for i, c in enumerate(s):
            indexMap[c].append(i)

        res = 0
        hit = 0
        while hit < len(t):
            hit = match(hit)
            res += 1
        return res


print(Solution().solve(s="dab", t="abbd"))
# If we concatenate a = "dab" three times, we can get "dabdabdab".
# And "abbd" is a subsequence of "dabdabdab".
# 3
