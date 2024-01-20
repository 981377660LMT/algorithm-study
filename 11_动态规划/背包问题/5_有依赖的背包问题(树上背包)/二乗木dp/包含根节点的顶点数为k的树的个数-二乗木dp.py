# 给定一棵n个顶点的树
# 求包含根节点0、顶点数为k的树的个数
# !结点编号为0~n-1
# https://snuke.hatenablog.com/entry/2019/01/15/211812

MOD = 998244353
n = int(input())
k = int(input())
edges = tuple(tuple(map(int, input().split())) for _ in range(n - 1))
adjList = [[] for _ in range(n)]
for u, v in edges:
    adjList[u].append(v)
    adjList[v].append(u)


# !dp[i][v]表示以i为根的子树中包含v个顶点的树的个数
# 1.n,k差不多大
# n<=3000
# k<=3000
# O(n^2)
# 咋一看这样的复杂度会是O(n^3)，
# 但如果每次枚举的范围都是儿子树的大小，可以证明这样的树型 dp的复杂度是 O(n^2)的。 (完全图)


def dfs1(cur: int, pre: int) -> None:
    subSize[cur] = 1
    dp[cur] = [0, 1]
    for next in adjList[cur]:
        if next == pre:
            continue
        dfs1(next, cur)
        merged = [0] * (subSize[cur] + subSize[next] + 1)  # 可以用卷积优化
        for i in range(1, subSize[cur] + 1):
            for j in range(1, subSize[next] + 1):
                merged[i + j] += dp[cur][i] * dp[next][j] % MOD
        subSize[cur] += subSize[next]
        dp[cur] = merged
    dp[cur][0] = 1


subSize = [0] * n
dp = [[] for _ in range(n)]
dfs1(0, -1)
print(dp[0][k])

# 2.k很小
# n<=1e5
# k<=500
# O(n*k)
# dpテーブルのサイズが K を超えたら K になるようにカット

# n个集合合并,一开始每个集合只有一个元素
# 合并两个集合的代价为min(|A|, K) * min(|B|, K)
# 合并成一个集合的最大代价为O(n*K)


def dfs2(cur: int, pre: int) -> None:
    global res
    subSize[cur] = 1
    dp[cur] = [0, 1]
    for next in adjList[cur]:
        if next == pre:
            continue
        dfs2(next, cur)
        merged = [0] * (subSize[cur] + subSize[next] + 1)
        for i in range(1, subSize[cur] + 1):
            for j in range(1, subSize[next] + 1):
                merged[i + j] += dp[cur][i] * dp[next][j] % MOD
        subSize[cur] += subSize[next]
        dp[cur] = merged
        if subSize[next] > k:  # 关键
            subSize[cur] = k
            dp[cur] = dp[cur][: k + 1]
    if subSize[cur] >= k:
        res += dp[cur][k]
        res %= MOD
    dp[cur][0] = 1


subSize = [0] * n
dp = [[] for _ in range(n)]
res = 0
dfs2(0, -1)
print(res)
