# D - Forbidden Difference (打家劫舍取模版本)
# https://atcoder.jp/contests/abc403/tasks/abc403_d
# 给定一个长度为N的数组以及一个非负整数D。
# 问至少删除多少个数字，才能够使得数组内不存在任何两个差值为D的元素。
#
# 差值为D => 模D相同

from collections import defaultdict
from itertools import groupby
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, D = map(int, input().split())
    A = list(map(int, input().split()))

    if D == 0:
        print(N - len(set(A)))
        exit(0)

    groups = defaultdict(list)
    for v in A:
        groups[v % D].append(v // D)
    for vs in groups.values():
        vs.sort()

    res = 0
    for group in groups.values():
        if len(group) == 0:
            continue

        # !dp: 最后一个元素不选还是选，可以保留的最多数字
        pre = -INF
        dp1, dp2 = 0, 0
        for v, c in groupby(group):
            cnt = len(list(c))
            if v == pre + 1:
                dp1, dp2 = max(dp1, dp2), dp1 + cnt
            else:
                dp1 = max(dp1, dp2)
                dp2 = dp1 + cnt
            pre = v
        res += max(dp1, dp2)

    print(N - res)
