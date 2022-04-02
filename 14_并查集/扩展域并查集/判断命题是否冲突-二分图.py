# n组命题，表示区间[a,b]里单词的个数是奇数还是偶数
# 问命题是否冲突
from collections import defaultdict


class Solution:
    def solve(self, lists):
        adjMap = defaultdict(list)
        for u, v, color in lists:
            adjMap[u].append([v + 1, color])
            adjMap[v + 1].append([u, color])

        colors = dict()

        def dfs(node, color):
            colors[node] = color
            for next, weight in adjMap[node]:
                if next not in colors:
                    if not dfs(next, color ^ weight):
                        return False
                elif colors[next] != color ^ weight:
                    return False

            return True

        for cur in adjMap:
            if cur not in colors and not dfs(cur, 0):
                return False
        return True


print(Solution().solve(lists=[[1, 5, 1], [6, 10, 0], [1, 10, 0]]))
# If there are an odd number of words from pages [1, 5] and an even number of words from pages [6, 10],
# that must mean there's an odd number from pages [1, 10]. So this is a contradiction.
