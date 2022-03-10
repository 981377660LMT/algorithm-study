from sys import setrecursionlimit


def main():
    MAX_N = 1505
    dp = [[0] * 3 for _ in range(MAX_N)]
    kid, w = [0] * MAX_N, [0] * MAX_N
    web = [[] for _ in range(MAX_N)]
    setrecursionlimit(10000)
    n = int(input())

    # 建图、保存权值以及找根结点
    for _ in range(n):
        info = list(map(int, input().split()))
        u, cost, m = info[:3]
        w[u] = cost
        for v in info[3:]:
            web[u].append(v)
            kid[v] = 1

    root = 1
    while kid[root]:
        root += 1

    # 状态机树形 DP，求最小花费
    # dp[u][0] 当前点 u 被父结点覆盖
    # dp[u][1] 当前点 u 被子结点覆盖
    # dp[u][2] 当前点 u 驻兵
    def dfs(root):
        # 当前点驻兵花费
        dp[root][2] = w[root]
        tot = 0
        for v in web[root]:
            dfs(v)
            # 当前点 u 被父结点覆盖，子结点 v 只能自力更生或者被 v 的子结点覆盖
            dp[root][0] += min(dp[v][1:])
            # 当前点 u 驻兵，子结点 v 选最小花费即可
            dp[root][2] += min(dp[v])
            # 预先统计子结点不被父结点覆盖的最小化花费之和
            tot += min(dp[v][1:])
        # 当前点 u 被子结点覆盖的最小花费 = tot + 当前子节点 v 驻兵 - 子节点 v 预先贡献给 tot 的值
        dp[root][1] = min([tot - min(dp[v][1:]) + dp[v][2] for v in web[root]] + [1e9])

    dfs(root)
    # 根结点没有父结点，因此在后两项中取最小值
    print(min(dp[root][1:]))


if __name__ == "__main__":
    main()
