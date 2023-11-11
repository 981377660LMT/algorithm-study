# 有n个人坐成一个桌，每个人面前都有一盆菜，
# 人和菜都是一个1-n的排列，如果一个菜距离自己的距离不超过1，那么这个人就会很高兴。
# 你作为服务员可以任意旋转桌子，使得高兴的人最多。
# !看转几次之后会产生贡献 再全部叠加

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
nums = list(map(int, input().split()))  # 每个人面前的料理编号
counter = [0] * (n + 5)
mp = {num: i for i, num in enumerate(nums)}
for i in range(n):
    for need in (i - 1, i, i + 1):  # 每个人都有三个可能
        pos = need % n
        index = mp[pos]
        dist = (i - index) % n
        counter[dist] += 1
print(max(counter))
