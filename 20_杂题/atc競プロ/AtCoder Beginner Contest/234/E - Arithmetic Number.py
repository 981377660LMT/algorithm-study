# !找到第一个大于等于 X 且满足要求的数字(每个数位等差,等差数)。
# 1<=x<=1e17

# !等差数列由首项和公差唯一确定
# !枚举第一位数是什么和公差是多少即可。时间复杂度：O(9∗19∗17)

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

x = int(input())

res = INF


for first in range(1, 10):
    for delta in range(-9, 10):
        arr = []
        cur = first
        while True:
            arr.append(cur)
            cand = int("".join(map(str, arr)))
            if cand >= x:
                res = min(res, cand)

            cur += delta
            if cur < 0 or cur >= 10 or len(arr) >= 20:
                break
print(res)
