from collections import defaultdict

# 最大权重的相交区间


class Solution:
    def solve(self, intervals):
        diff = defaultdict(int)
        for s, e, w in intervals:
            diff[s] += w
            diff[e + 1] -= w

        curSum = 0
        maxSum = 0
        events = []
        for time in sorted(diff):
            curSum += diff[time]
            maxSum = max(maxSum, curSum)
            events.append([time, curSum])

        res = []
        for i, (s, w) in enumerate(events):
            if w != maxSum:
                continue
            if res and res[-1][-1] + 1 == s:
                res[-1][-1] = events[i + 1][0] - 1
            elif i + 1 < len(events):
                res.append([s, events[i + 1][0] - 1])
        return res


print(Solution().solve(intervals=[[1, 3, 1], [2, 6, 1], [5, 7, 1]]))

# [
#     [2, 3],
#     [5, 6]
# ]
# During [[2, 3], [5, 6]] weight is 2 which is the max possible.
# Return the list of intervals that have the highest weight, sorted in ascending order
