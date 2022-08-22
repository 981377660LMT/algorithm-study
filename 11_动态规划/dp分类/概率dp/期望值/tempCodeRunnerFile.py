
# nums = list(map(float, input().split()))
# counter = Counter(nums)


# memo = [-1.0] * (n + 1) * (n + 1) * (n + 1)


# def dfs(remain1: int, remain2: int, remain3: int) -> float:
#     """counter保存状态

#     这道题的加强版 九坤t4-筹码游戏-组合.py
#     """
#     if remain1 == remain2 == remain3 == 0:
#         return 0
#     hash_ = remain1 * n * n + remain2 * n + remain3
#     if memo[hash_] != -1:
#         return memo[hash_]

#     div = remain1 + remain2 + remain3
#     res = n / div
#     p1, p2, p3 = remain1 / div, remain2 / div, remain3 / div
#     if remain1:
#         res += p1 * dfs(remain1 - 1, remain2, remain3)
#     if remain2:
#         res += p2 * dfs(remain1 + 1, remain2 - 1, remain3)
#     if remain3:
#         res += p3 * dfs(remain1, remain2 + 1, remain3 - 1)
#     memo[hash_] = res
#     return res


# print(dfs(counter[1], counter[2], counter[3]))