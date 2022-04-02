# 那些顾客可能是第k(0-indexed)个访问网站的顾客
class Solution:
    def solve(self, requests, k):
        # calculate the earliest and latest visit times for user k.
        # return all requests that overlap with this interval.
        starts = sorted(s for s, e in requests)[k]
        ends = sorted(e for s, e in requests)[k]
        return [i for i, (s, e) in enumerate(requests) if not (ends < s or e < starts)]


print(Solution().solve(requests=[[3, 4], [1, 3], [4, 4]], k=1))
