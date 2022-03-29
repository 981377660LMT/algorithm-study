# 敏捷开发
from collections import defaultdict


class Solution:
    def solve(self, intervals, types):
        """Return a sorted merged list of lists where each element contains [start, end, num_types]"""
        dic = defaultdict(int)
        for left, right in intervals:
            dic[left] += 1
            dic[right] -= 1

        res = []
        pre = -1
        for cur in sorted(dic):
            if dic[pre] > 0:
                res.append([pre, cur, dic[pre]])
            dic[cur] += dic[pre]
            pre = cur
        return res


print(
    Solution().solve(
        intervals=[[0, 3], [5, 7], [0, 7]], types=["hacker news", "reddit", "scrum meeting"]
    )
)
