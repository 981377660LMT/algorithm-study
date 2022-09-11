# 环上的差分

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
nums = list(map(int, input().split()))
diff = [0] * (n + 5)
mp = {num: i for i, num in enumerate(nums)}
for i in range(n):
    for cand in (i - 1, i, i + 1):
        val = cand % n
        index = mp[val]
        dist = (i - index) % n
        diff[dist] += 1
print(max(diff))
