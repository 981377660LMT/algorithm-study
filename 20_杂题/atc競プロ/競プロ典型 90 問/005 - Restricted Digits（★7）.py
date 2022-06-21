# K 種類の数字 c1​,c2​,…,cK​ のみを使うことで作れる N 桁の正の整数のうち、
# B の倍数は何個あるか、1000000007 で割ったあまりを答えてください。

# 1<=N<=1e18
# 2<=B<=1e3
# 1<=c1<...<ck<=9

# 倍数构造的方案数
# !数位dp dp[i+1][(mod*10+select)%B]+=dp[i][mod]

from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e9))

MOD = int(1e9 + 7)

N, B, K = map(int, input().split())
nums = list(map(int, input().split()))

# 显然这个1e18是不可能线性遍历的 必须倍增处理 O(B^2logN)
# https://drken1215.hatenablog.com/entry/2021/10/08/231200

# 倍增法的关键:
# !dp[i+1][j]=dp[i][dp[i][j]]
