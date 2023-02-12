# https://atcoder.jp/contests/dwango2016-prelims/tasks/dwango2016qual_e
# 公司位于0 家位于L
# 一个人从公司走回家,不能回头、不能停下来
# 途中正好在放烟花
# 烟花一共n发,第i发烟花在ti时刻放出,烟花的位置为pi
# !这个人想让自己的位置和烟花的位置的距离的和最小
# 求出最小的和

# 给定数列nums1
# !求一个严格单调递增的数列nums2使得∑(nums1[i]-nums2[i])最小
# !dp[i][x] = min(dp[i-1][y] + abs(x-pi)) , y<=x

import sys
from typing import List, Tuple
from SlopeTrick import SlopeTrick

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(home: int, fireworks: List[Tuple[int, int]]) -> int:
    # dp[i][x] = min(dp[i-1][y] + abs(x-nums1[i])) , y<=x
    fireworks.sort(key=lambda x: (x[0], -x[1]))  # !按时间排序,时间相同时按照位置降序
    S = SlopeTrick()
    for _, pos in fireworks:
        S.clear_right()
        S.add_abs(pos)
    return S.query()[0]


if __name__ == "__main__":
    n, home = map(int, input().split())
    fireworks = []  # (time,pos)
    for _ in range(n):
        time, pos = map(int, input().split())
        fireworks.append((time, pos))

    print(solve(home, fireworks))
