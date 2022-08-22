# あなたは以下の操作を M 回以下行うことができます
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m, k = map(int, input().split())
nums = list(map(int, input().split()))


# !贪心 能不能加到某个二进制位
# 二分答案

# 最后k个数最好都一样
counter = [0] * 32
for num in nums:
    for i in range(32):
        if num & (1 << i):
            counter[~i] += 1
print(counter)
