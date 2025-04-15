# abc400-F-环染色的最小代价-环形区间dp
# https://atcoder.jp/contests/abc400/submissions/64615934
#
# 给你一个环，每次可以选出一个区间，然后给那个区间染成任意颜色 c，代价是区间长度加 X[c]，求最小代价使得最终颜色序列为目标序列C.
# n<=400
#
# dp[l][r]表示从 l 到 r 的区间的最小代价.
# !翻倍数组，我们只要通过区间DP算出倍长后数组中所有长度为 n 的区间的代价，就可以求出将整块蛋糕染色的总代价了.


from functools import lru_cache


INF = int(1e18)

if __name__ == "__main__":
    n = int(input())
    targets = list(map(int, input().split()))
    targets = [x - 1 for x in targets]
    costs = list(map(int, input().split()))

    targets2 = targets + targets

    @lru_cache(None)
    def dfs(left: int, right: int) -> int:
        """闭区间[left,right]这一段染色的最小代价."""
        if left > right:
            return 0
        if left == right:
            return 1 + costs[targets2[left]]

        res = dfs(left, right - 1) + 1 + costs[targets2[right]]  # !染成右边的颜色

        # !染这一段
        for mid in range(left, right):
            if targets2[mid] == targets2[right]:
                res = min(
                    res,
                    dfs(left, mid) + dfs(mid + 1, right - 1) + (right - mid),
                )
        return res

    res = INF
    for i in range(n):
        res = min(res, dfs(i, i + n - 1))
    dfs.cache_clear()
    print(res)
