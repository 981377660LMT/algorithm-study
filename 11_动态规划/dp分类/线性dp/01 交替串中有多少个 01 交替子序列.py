# 求一个长为 n(<=1e6) 的 01 交替串中有多少个 01 交替子序列。
# 对结果模 1e9+7。

# 输入 n=3
# 输出 6
# 解释 交替串 101 有如下 01 交替子序列（x 表示不选）
# 1xx
# x0x
# xx1
# 10x
# x01
# 101
# 你也可以用 010 当作交替串，算出来的结果仍然是 6。

# 输入 n=4
# 输出 11

# !每个位置记录结尾为0，结尾为1的符合条件的个数


MOD = int(1e9 + 7)
n = int(input())
# dp = [0] * (n + 1)
# dp[1] = 1
# for i in range(2, n + 1):
#     dp[i] = dp[i - 1] + dp[i - 2] + 2
#     dp[i] %= MOD
# print(dp[n])

endswith = [0, 0]
for i in range(n):
    num = 0 if i % 2 == 0 else 1
    endswith[num] += endswith[num ^ 1] + 1
    endswith[num] %= MOD
print(sum(endswith) % MOD)

