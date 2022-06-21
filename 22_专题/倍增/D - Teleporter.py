# 町が N 個ある。町 i から町 Ai に移動することを K 回繰り返す。
# 町 1 から始めた時、最終的にどの町にたどり着くか？
# N<=2e5 K<=1e18


from math import floor, log2
import sys

input = sys.stdin.readline

n, k = map(int, input().split())

nexts = [int(v) - 1 for v in input().split()]  # 町 i(1≤i≤N) のテレポーターの転送先は町 A iです


# 解法一 哈希表+查找循环节(类似于 n天后的牢房)


# !解法二 倍增
maxJ = floor(log2(k)) + 1
# doubling[k][i] : 町 i から 2^k 先の町はどこか
# 最后计算时将k二进制分解即可

# 1. 初始化一步
dp = [[0] * (n + 1) for _ in range(maxJ + 1)]
for i in range(n):
    dp[0][i] = nexts[i]

# 2. 倍增
for j in range(maxJ):  # j+1<=maxJ
    for i in range(n):
        dp[j + 1][i] = dp[j][dp[j][i]]

# 3.二进制分解计算
res = 0
for bit in range(maxJ + 1):
    if (k >> bit) & 1:
        res = dp[bit][res]

print(res + 1)

