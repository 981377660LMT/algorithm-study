from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
nums = list(map(int, input().split()))


delta = [0] * (n + 5)
mp = {num: i for i, num in enumerate(nums)}
# 不移动时的值
res = 0

for i in range(n):
    index = mp[i]
    dist = min((i - index) % n, (index - i) % n)
    res += dist
mid = n // 2
for i in range(n):
    index = mp[i]
    if index >= i:
        count = index - i
        delta[1] -= 1
        delta[count + 1] += 1
    else:
        count = i - index
        delta[1] -= 1
        delta[count + 1] += 1
        delta[n - count + 1] -= 1
        delta[n + 1] += 1

diff = list(accumulate(delta))
print(diff)
print(res)

# 差分
