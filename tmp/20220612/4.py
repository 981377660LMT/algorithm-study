from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 2 <= ideas.length <= 5 * 104
# 1 <= ideas[i].length <= 10
# ideas[i] 由小写英文字母组成
# !返回 不同 且有效的公司名字的数目。
# !交换前缀不就等于交换后缀吗


class Solution:
    def distinctNames(self, ideas: List[str]) -> int:
        n = len(ideas)

        # 不合法的
        adjMap = defaultdict(set)
        startCounter = defaultdict(int)
        for word in ideas:
            # 长度为1怎么办
            adjMap[word[1:]].add(word[0])
            startCounter[word[0]] += 1
        # defaultdict(<class 'set'>, {'offee': {'c', 't'}, 'onuts': {'d'}, 'ime': {'t'}})
        # defaultdict(<class 'int'>, {'c': 1, 'd': 1, 't': 2})
        res = n * (n - 1)
        for chars in adjMap.values():
            allSum = sum(startCounter[char] for char in chars)
            for char in chars:
                res -= (allSum - startCounter[char]) * startCounter[char]
        return res


print(Solution().distinctNames(ideas=["coffee", "donuts", "time", "toffee"]))
print(Solution().distinctNames(ideas=["lack", "back"]))
print(Solution().distinctNames(["aaa", "baa", "caa", "bbb", "cbb", "dbb"]))
# class Solution:
#     def distinctNames(self, ideas: List[str]) -> int:
#         d = defaultdict(set)
#         for name in ideas:
#             d[name[0]].add(name[1:])
#         res = 0
#         # print(d)
#         for a, b in product(d, repeat=2):
#             if a == b: continue
#             na, nb = len(d[a]), len(d[b])
#             nc = len(d[a] & d[b])
#             res += (na - nc) * (nb - nc)
#             # print(na, nb, nc, a, d[a], b, d[b], d[a] & d[b], (na - nc) * nb)

#         return res

# class Solution:
#     def calculateTax(self, brackets: List[List[int]], income: int) -> float:
#         res = cur = 0
#         for m, p in brackets:
#             if income == 0: break
#             r = min(m - cur, income)
#             res += r * p / 100
#             cur += r
#             income -= r
#         return res


# class Solution:
#     def minPathCost(self, grid: List[List[int]], moveCost: List[List[int]]) -> int:
#         m, n = len(grid), len(grid[0])
#         dp = [grid[0][j] for j in range(n)]
#         for i in range(1, m):
#             ndp = [1234567890123] * n
#             for j in range(n):
#                 ndp[j] = grid[i][j] + min(dp[k] + moveCost[grid[i - 1][k]][j] for k in range(n))
#             dp = ndp
#         return min(dp)
