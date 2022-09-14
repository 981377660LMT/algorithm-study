# 给定N和M，求满足以下条件的正整数序列的个数:
# ·长度为N，每个元素∈[1, M]。
# ·其LIS长度恰好为3。
# 3≤N<1000，3< M≤10。

# !dp[i][a][b][c] 表示前i个数 长度为1 2 3 的LIS结尾为a b c时 序列的个数
# !O(n*m^4) 1e8计算
# LIS长为3的子序列个数


MOD = 998244353
INF = int(4e18)

n, m = map(int, input().split())

dp = {(INF, INF, INF): 1}  # !这个和三维数组差不多快
for _ in range(n):
    ndp = dict()
    for pre, count in dp.items():
        for cur in range(1, m + 1):
            next_ = ()
            if cur <= pre[0]:
                next_ = (cur, pre[1], pre[2])
            elif cur <= pre[1]:
                next_ = (pre[0], cur, pre[2])
            elif cur <= pre[2]:
                next_ = (pre[0], pre[1], cur)
            else:
                continue
            ndp[next_] = (ndp.get(next_, 0) + count) % MOD
    dp = ndp


res = 0
for k, v in dp.items():
    if k[2] != INF:
        res += v
        res %= MOD
print(res)
