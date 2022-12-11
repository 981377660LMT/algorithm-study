# 给你一个由n个点组成的环，1 <=> 2 <=> 3 <=> ... <=> n <=> 1
# 在给一个0号节点,与所有的节点都有边连接。
# 所有的边都有边权.
# !现在你需要删除一些边，使得这个图是一个二分图，问删除的权值最小是多少。
# n<=2e5

# !二分图+环上dp分类讨论
# 由二分图对称性 设0结点的颜色为0
# 对于点i 如果
# !1. i和i-1颜色相同 则需要删除i-1到i的边
# !2. i和0颜色相同 则需要删除0到i的边


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":

    n = int(input())
    score1 = list(map(int, input().split()))  # !i到0的边权
    score2 = list(map(int, input().split()))  # !i到(i+1)的边权

    # # !记忆化搜索 python3.8
    # @lru_cache(None)
    # def dfs(index: int, pre: int, root: int) -> int:
    #     """当前在index 前一个颜色为pre 第一个点的颜色为root"""
    #     if index == n + 1:
    #         return score2[-1] if pre == root else 0

    #     # 当前点的颜色为0
    #     cand1 = dfs(index + 1, 0, root) + score1[index - 1] + (score2[index - 2] if pre == 0 else 0)

    #     # 当前点的颜色为1
    #     cand2 = dfs(index + 1, 1, root) + (score2[index - 2] if pre == 1 else 0)

    #     return min(cand1, cand2)

    # # 1颜色为0
    # res0 = dfs(2, 0, 0) + score1[0]
    # # 1颜色为1
    # res1 = dfs(2, 1, 1)
    # print(min(res0, res1))
    ##############################################################################
    # !dp[i][j][k]表示i号点时 前一个点的颜色为j 第一个点的颜色为k时的最小值
    dp = [[[INF, INF] for _ in range(2)] for _ in range(n + 1)]
    dp[1][0][0] = score1[0]
    dp[1][1][1] = 0
    for i in range(2, n + 1):
        for pre in range(2):
            for cur in range(2):
                for root in range(2):
                    dp[i][cur][root] = min(
                        dp[i][cur][root],
                        dp[i - 1][pre][root]
                        + (score2[i - 2] if pre == cur else 0)  # 前一个点的颜色和当前点的颜色相同
                        + (score1[i - 1] if cur == 0 else 0),  # 当前点的颜色为0
                    )

    res = INF
    for pre in range(2):
        for root in range(2):
            res = min(res, dp[n][pre][root] + (score2[-1] if pre == root else 0))  # !最后一个点与root是否同色
    print(res)
