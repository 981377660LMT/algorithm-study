# K 種類の数字 c1​,c2​,…,cK​ のみを使うことで作れる N 桁の正の整数のうち、
# B の倍数は何個あるか、1000000007 で割ったあまりを答えてください。

# 1<=N<=1e18
# 2<=B<=1e3
# 1<=c1<...<ck<=9

# 倍数构造的方案数 i桁決めてmodBでjになるものの数
# !数位dp dp[i+1][(mod*10+select)%B]+=dp[i][mod]
# dp[i]→dp[i+1] とすると地球爆発　jとｃだけで変わる
# 半分に分けて最後に合体させる


# !数位dp+倍增
from math import floor, log2
import sys

sys.setrecursionlimit(int(1e9))

input = lambda: sys.stdin.readline()

MOD = int(1e9 + 7)

n, b, k = map(int, input().split())  # 位数 倍数 k种数字
canUse = list(map(int, input().split()))

# 显然这个1e18是不可能线性遍历的 必须倍增处理 O(B^2logN)
# !把数字分为前半和后半 dp[i + j][(p * tj + q) % B] += dp[i][p] * dp[j][q]
# https://drken1215.hatenablog.com/entry/2021/10/08/231200
# https://atcoder.jp/contests/typical90/submissions/32340943

# maxJ = floor(log2(k)) + 1
maxJ = 60
pow10 = [pow(10, 2 ** i, b) for i in range(maxJ + 1)]

# 1. 初始化
dp = [[0] * b for _ in range(maxJ + 1)]
for num in canUse:
    dp[0][num % b] += 1

# 2. 分开之后倍增
for j in range(maxJ):
    for i1 in range(b):
        for i2 in range(b):
            dp[j + 1][(pow10[j] * i1 + i2) % b] += dp[j][i1] * dp[j][i2]
            dp[j + 1][(pow10[j] * i1 + i2) % b] %= MOD


res = [0] * b
res[0] = 1
for j in range(maxJ + 1):
    if (n >> j) & 1:
        nRes = [0] * b
        for i1 in range(b):
            for i2 in range(b):
                nRes[(pow10[j] * i1 + i2) % b] += res[i1] * dp[j][i2]
                nRes[(pow10[j] * i1 + i2) % b] %= MOD
        res = nRes

print(res[0])
