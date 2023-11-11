# 给定n个数对(A1,B1), (A2, B2),.. . , (An, Bn)
# !对于一段连续区间[l, r] (1<=l<=r<=m)，
# !如果对于所有n 个数对: (Ai,Bi)中至少有一个数字在区间[l, r]中，
# 那么说明这个[l, r]区间是个好区间。
# !问:长度为1,2,3,.. . , m的好区间分别有多少个?
# n <= 2e5,2≤m≤2e5,1 < Ai,Bi≤m


# !1.注意到滑窗包含数字个数的单调性
# !2.如果我们找到一个最小的满足条件的子区间 [l,r] , 那就不用去枚举包含该子区间的大区间了。
# !3.差分数组做区间更新

# !因此只需要固定左端点看右端点的边界即可
from collections import defaultdict
from itertools import accumulate
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    pair = [tuple(map(int, input().split())) for _ in range(n)]
    indexMap = defaultdict(list)
    for i, (a, b) in enumerate(pair):
        indexMap[a].append(i)
        indexMap[b].append(i)

    diff = [0] * (m + 10)

    right, counter = 1, defaultdict(int)
    for left in range(1, m + 1):
        while len(counter) < n and right <= m:
            for index in indexMap[right]:
                counter[index] += 1
            right += 1

        # 此时的r实际为r-1
        if len(counter) == n:
            lower = (right - 1) - left + 1
            upper = m - left + 1
            diff[lower] += 1
            diff[upper + 1] -= 1

        for index in indexMap[left]:
            counter[index] -= 1
            if counter[index] == 0:
                del counter[index]

    diff = list(accumulate(diff))
    for i in range(1, m + 1):
        print(diff[i])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
