# One Interval
# 增添一个最短的区间使得所有区间都连续


class Solution:
    def solve(self, intervals):
        if len(intervals) == 1:
            return 0

        intervals.sort()

        res = [intervals[0]]
        for i in range(1, len(intervals)):
            if res[-1][1] >= intervals[i][0]:
                res[-1][1] = max(res[-1][1], intervals[i][1])
            else:
                res.append(intervals[i])

        if len(res) == 1:
            return 0

        return res[-1][0] - res[0][1]  # 注意这里
