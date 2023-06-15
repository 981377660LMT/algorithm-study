# https://atcoder.jp/contests/abc268/tasks/abc268_e
# 中国餐厅（三星版本）
# 有n个人坐成一桌，定义一个人的沮丧程度为他喜欢的菜距离他的最小值(距离)，
# 请随意旋转桌子，最小化`所有人的沮丧程度之和`。

# !环上的差分 一次函数叠加
# !一个人的沮丧程度可以画成一个很多一次函数组成的函数
# 折线在相加时k和b其实是互相独立的，
# !因此我们对一段区间加上一条直接可以分别加k和b，
# 区间加变成差分单点加，加完后做一遍前缀和，On扫一遍求最小值
# 累加所有的折线图，然后找到新的图的最小值。

# https://zhuanlan.zhihu.com/p/563283585

from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
nums = list(map(int, input().split()))


def update(b0: int) -> None:
    def add(left: int, right: int, deltaK: int, deltaB: int) -> None:
        if left > right:
            return
        diffK[left] += deltaK
        diffK[right + 1] -= deltaK
        diffB[left] += deltaB
        diffB[right + 1] -= deltaB

    mid = n // 2
    if b0 < mid:
        add(0, b0 - 1, -1, b0)
        add(b0, b0 + mid, 1, -b0)
        add(b0 + mid + 1, n - 1, -1, b0 + n)
    else:
        add(0, b0 - mid - 1, 1, n - b0)
        add(b0 - mid, b0 - 1, -1, b0)
        add(b0, n - 1, 1, -b0)


diffK = [0] * (n + 5)  # 斜率的差分
diffB = [0] * (n + 5)  # 截距的差分
mp = {num: i for i, num in enumerate(nums)}
for i in range(n):
    pos = mp[i]
    dist = (i - pos) % n
    update(dist)

diffK = list(accumulate(diffK))
diffB = list(accumulate(diffB))
res = INF
for i in range(n):
    cand = diffK[i] * i + diffB[i]
    res = min(res, cand)
print(res)
