# https://atcoder.jp/contests/tdpc/editorial/752

# 构造一个小写字母构成的字符串，使小写字母 i 出现了 freqi​ 次
# !且字符串相邻的两个字符不能相同
# 求字符串的个数
# freqi<=10 (即长度<=260)

# !dp[i][a] 表示前i种字符 相邻字符出现连续的次数为a 答案为dp[-1][0]
# 隣り合う箇所 a を決めると隣り合わない箇所の数も一意に決まるのでこれを b とします
# !插入新的字符 可以插入连续槽a和不连续槽b
# !将这f个新字符分为g组(1<=g<=f)插入 g组里选择j个连续槽a (g-j)个不连续槽b
# !分别为 C(f-1,g-1) C(a,j) C(b,g-j)
# !新的连续槽个数为 a+f-g-j

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")

#########################################################
from itertools import accumulate

MOD = int(1e9 + 7)

fac = [1]
ifac = [1]
for i in range(1, int(1e4)):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


##################################################################
freq = [num for num in (map(int, input().split())) if num > 0]
preSum = [0] + list(accumulate(freq))
n = len(freq)

fsum = sum(freq)
dp = [0] * (fsum + 5)
dp[0] = 1

for ci, cc in enumerate(freq):
    ndp = [0] * (fsum + 5)
    for adj, preRes in enumerate(dp):  # 相邻插入槽adj
        if preRes == 0:
            continue
        notAdj = (preSum[ci] + 1) - adj  # !注意不相邻插入槽有preSum[ci]+1个, 3个字符分隔出4个插入槽
        for group in range(1, cc + 1):
            count1 = C(cc - 1, group - 1)
            for slot1 in range(min(group, adj) + 1):
                count2 = C(adj, slot1)
                count3 = C(notAdj, group - slot1)
                ndp[adj + cc - group - slot1] += preRes * count1 % MOD * count2 % MOD * count3
                ndp[adj + cc - group - slot1] %= MOD
    dp = ndp

print(dp[0] % MOD)
