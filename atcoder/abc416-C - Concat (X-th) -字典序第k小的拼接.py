# abc416-C - Concat (X-th) -字典序第k小的拼接
# https://atcoder.jp/contests/abc416/tasks/abc416_c
# 给定 N 个字符串 S₁,…,Sₙ。考虑所有长度为 K 且元素在 1…N 范围内的数列 (A₁,…,Aₖ)，定义
# f(A₁,…,Aₖ) = S_{A₁} + S_{A₂} + … + S_{Aₖ} （字符串拼接）。
# 将所有 N^K 个 f(...) 按字典序排序，求第 X 小的字符串。

N, K, X = map(int, input().split())
words = [input().strip() for _ in range(N)]

dp = [""]
for _ in range(K):
    ndp = []
    for pre in dp:
        for word in words:
            ndp.append(pre + word)
    dp = ndp

dp.sort()
print(dp[X - 1])

# 加强版：
# https://atcoder.jp/contests/abc416/tasks/abc416_g
# !1<=N,K<=1e5
# !每个单词长度不超过10
