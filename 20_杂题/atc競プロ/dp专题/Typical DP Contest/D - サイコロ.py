# 掷骰子n次 求掷出的点数之积为d的倍数的概率
# n<=100
# d<=1e18

# !1. 质因子必须为2/3/5 因此可能的倍数为 O(logd)**3 个
# !2. dp[index][mul] 来dp `注意mul需要模d` 取模之后不存在的素因子还是会不存在
from collections import defaultdict
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def check(num: int) -> bool:
    for f in (2, 3, 5):
        while num % f == 0:
            num //= f
    if num != 1:
        return False
    return True


n, d = map(int, input().split())
if not check(d):
    print(0)
    exit()

dp = defaultdict(int, {i: 1 / 6 for i in range(1, 7)})
for i in range(1, n):
    ndp = defaultdict(lambda: 0.0)
    for pre, value in dp.items():
        for cur in range(1, 7):
            ndp[pre * cur % d] += value / 6
    dp = ndp


print(dp[0])
