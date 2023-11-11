# 这里有N(N ≤300)个人，每个人都参加了两场比赛，第i的人的排名分别为a[i],b[i]。
# 现在需要选择K个人去参加比赛，这K人需要满足以下一个条件。

# 如果X被选择了而且Y没有被选择，那么a[x] > a[y] 并且 b[x] > b[Y]
# 也就是说，如果 y 的两场的排名都比 x 小，不能出现选择 x 而不选择 y 的情况
# 求选择方案数 mod 998244353

# !1. 二维偏序需要对某个维度排序 维护另一个维度
# !2. 1<=k<=n<=300 暗示dp O(n^3)
# 二维偏序 + dfs(index,remain) 模型

# 按照第一次排名升序遍历选还是不选dp dp过程中记录之前没有选择的人中第二次的最小排名
# 如果当前人的第二次排名比之前的最小排名大，那么就不能选择这个人参赛

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 注意dfs要用数组+哈希存状态快一些
def dfs(index: int, remain: int, preMin: int) -> int:
    if remain < 0:
        return 0
    if index == n:
        return 1 if remain == 0 else 0
    hash = index * h2 + remain * h1 + preMin
    if memo[hash] != -1:
        return memo[hash]

    res = dfs(index + 1, remain, min(preMin, people[index][1]))  # jump
    if remain > 0 and people[index][1] < preMin:
        res += dfs(index + 1, remain - 1, preMin)  # not jump
    res %= MOD
    memo[hash] = res
    return res


n, k = map(int, input().split())
rank1 = list(map(int, input().split()))
rank2 = list(map(int, input().split()))
people = sorted(zip(rank1, rank2), key=lambda x: x[0])
h3, h2, h1 = (n + 5) * (n + 5) * (n + 5), (n + 5) * (n + 5), (n + 5)
memo = [-1] * h3
print(dfs(0, k, n + 1))
