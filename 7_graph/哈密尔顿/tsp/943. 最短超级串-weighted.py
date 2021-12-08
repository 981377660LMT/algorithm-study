from typing import List
from collections import deque

# Travelling salesman problem (TSP) problem
# Time complexity: O(n^2 * 2^n)
# 对于一个字符串s和一个字符串t 如果s和t有v个字符重合 我们就连一条权值为w[s][t]=v的边，最终就可以得到一个完全图，那问题就转化为在这个完全图上求一条最长的路径


# 找到以 words 中每个字符串作为子字符串的最短字符串。如果有多个有效最短字符串满足题目条件，返回其中 任意一个 即可。

# 1 <= words.length <= 12
# 1 <= words[i].length <= 20
# words 中的所有字符串 互不相同


# summary:
# 这道题就转换成了，在一个图中，从某个点出发将所有点恰好遍历一遍，使得最后路过的路径长度最长。


class Solution:
    def shortestSuperstring(self, words: List[str]) -> str:
        def getWeight(s1: str, s2: str):
            for i in range(len(s1)):
                if s2.startswith(s1[i:]):
                    return len(s1) - i
            return 0

        def pathToString(path: List[int]) -> str:
            res = words[path[0]]
            for i in range(1, len(path)):
                commonLength = weight[path[i - 1]][path[i]]
                res += words[path[i]][commonLength:]
            return res

        n = len(words)
        target = (1 << n) - 1
        dist = [[0] * (1 << n) for _ in range(n)]
        weight = [[0] * n for _ in range(n)]
        for i in range(n):
            for j in range(i + 1, n):
                weight[i][j] = getWeight(words[i], words[j])
                weight[j][i] = getWeight(words[j], words[i])
        # bfs from each point
        # cur,state,path
        queue = deque([(i, 1 << i, [i], 0) for i in range(n)])

        bestPath = list(range(n))
        maxCost = -0x7FFFFFFF

        while queue:
            cur, state, path, cost = queue.popleft()

            # compare cost when finish
            if state == target:
                if cost > maxCost:
                    maxCost = cost
                    bestPath = path

            # visited
            if cost < dist[cur][state]:
                continue

            for next in range(n):
                if state & (1 << next):
                    continue
                nextState = state | (1 << next)
                if dist[cur][state] + weight[cur][next] > dist[next][nextState]:
                    dist[next][nextState] = dist[cur][state] + weight[cur][next]
                    queue.append((next, nextState, path + [next], dist[next][nextState]))
        # print(maxCost, bestPath)
        return pathToString(bestPath)


print(Solution().shortestSuperstring(words=["alex", "loves", "leetcode"]))
# 输出："alexlovesleetcode"
# 解释："alex"，"loves"，"leetcode" 的所有排列都会被接受。
print(Solution().shortestSuperstring(words=["catg", "ctaagt", "gcta", "ttca", "atgcatc"]))
