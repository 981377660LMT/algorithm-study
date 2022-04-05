from collections import defaultdict


# 二分+二分图
class Solution:
    def solve(self, n, views):
        if n <= 2:
            return 0

        views.sort(key=lambda x: x[2])

        # 排序后二分索引
        left, right = 0, len(views) - 1
        res = views[-1][2]
        while left <= right:
            mid = (left + right) >> 1
            cand = views[mid][2]
            if self.check(views, mid):
                res = min(res, cand)
                right = mid - 1
            else:
                left = mid + 1

        return res

    def check(self, views, start):
        adjMap = defaultdict(list)
        for i in range(start + 1, len(views)):
            a, b, _ = views[i]
            adjMap[a].append(b)
            adjMap[b].append(a)
        return self.isBipartite(adjMap, len(views))

    def isBipartite(self, graph, n) -> bool:
        def dfs(cur, c) -> bool:
            colors[cur] = c
            for neighbor in graph[cur]:
                if colors[neighbor] == -1:
                    if not dfs(neighbor, c ^ 1):
                        return False
                elif colors[neighbor] == c:
                    return False
            return True

        colors = [-1] * (n + 1)
        for i in range(n + 1):
            if colors[i] == -1 and not dfs(i, 0):
                return False
        return True


print(Solution().solve(n=3, views=[[0, 1, 100], [0, 2, 900], [1, 2, 500]]))
# 体育场需要有两个看台
# a和b比赛时 有c 个人看

# We can put teams 0 and 1 in one stadium and team 2 in the other stadium.
# This means team 2 doesn't have to compete with any other team.
# Then, for the first stadium we need capacity of 100 to host the game between teams 0 and 1.
