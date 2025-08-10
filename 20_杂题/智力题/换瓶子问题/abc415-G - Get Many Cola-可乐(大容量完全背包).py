# abc414-G - Get Many Cola-可乐(超大容量完全背包)
# https://atcoder.jp/contests/abc415/tasks/abc415_g
# https://zhuanlan.zhihu.com/p/1930037283559576140
#
# 有一个神奇的可乐店，只能用空瓶换新可乐。
# 初始有 N 瓶可乐，0 空瓶。每次可以：
#
# 喝掉一瓶可乐，可乐-1，空瓶+1。
# 选择某个 i（1≤i≤M），用 Ai 个空瓶换 Bi 瓶可乐（需空瓶≥Ai，且 Bi<Ai）。 问最多能喝多少瓶可乐。
#
# 1<=N<=1e15
# 1<=Bi<Ai<=300
#
# !总喝数 = N + max ∑Bi s.t. ∑(Ai−Bi) ≤ N，是一个容量为 N 的完全背包，物品数 M，物品 i 的“重量”Ci=Ai−Bi，价值Vi=Bi。
# !O(maxW**3)
#
# !大范围贪心(性价比)，小范围背包

N, M = map(int, input().split())
A, B = [0] * M, [0] * M
for i in range(M):
    A[i], B[i] = map(int, input().split())

maxA = max(A)
upper = maxA * (maxA + 1)
minLost = list(range(maxA + 1))  # minLost[c] -> 使用c个空瓶换可乐时，会损耗多少瓶子
for a, b in zip(A, B):
    minLost[a] = min(minLost[a], a - b)

dp = [0] * (upper + 1)  # dp[i] -> 使用i个空瓶时，能喝的最多可乐数
for i in range(upper + 1):
    for j in range(maxA + 1):
        if i >= j:
            dp[i] = max(dp[i], dp[i - minLost[j]] + j)

res = -int(1e18)
for i in range(1, maxA + 1):
    # 先一直用i操作，直到瓶子个数 <= upper，即大范围贪心
    if N < upper:
        ti = 0
    else:
        ti = (N - upper + minLost[i] - 1) // minLost[i]
    res = max(res, dp[N - ti * minLost[i]] + ti * i)

print(res)
