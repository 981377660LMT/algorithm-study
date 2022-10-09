"""
给定一个数组nums
每个回合对nums[i]增加(i+1)
求第1到第m个回合的数组的非负整数mex
n,m<=2e5 -1e9<=nums[i]<=1e9

!注意到只有[0,n]里的数才是可能的mex
对nums[0] 有n个数 nums[1] 有n/2个数 nums[2] 有n/3个数...
!最多保存nlogn个数
"""

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 以下の操作を M 回行ってください。
# 各 i (1≤i≤N) について、 A iに i を加算する。その後 A に含まれない最小の非負整数を求める。
# i (1≤i≤M) 行目には i 回目の操作後に A に含まれない最小の非負整数を出力せよ。
if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))

    def cal1(i: int) -> int:
        """最开始在第几个回合时>=0"""
        first = nums[i]
        left, right = 0, m
        while left <= right:
            mid = (left + right) // 2
            last = first + mid * (i + 1)
            if last >= 0:
                right = mid - 1
            else:
                left = mid + 1
        return left

    def cal2(i: int) -> int:
        """最后在第几个回合时<=n"""
        first = nums[i]
        left, right = 0, m
        while left <= right:
            mid = (left + right) // 2
            last = first + mid * (i + 1)
            if last <= n:
                left = mid + 1
            else:
                right = mid - 1
        return right

    # 小于0的和超过n的数不用管
    turn = [set() for _ in range(m + 10)]

    # 对每个数计算范围在[0,n]时,处于哪些回合
    # 二分
    for i in range(n):
        first = nums[i]
        round1, round2 = cal1(i), cal2(i)
        for round in range(round1, round2 + 1):
            cur = first + round * (i + 1)
            turn[round].add(cur)

    for t in range(1, m + 1):
        mex = 0
        while mex in turn[t]:
            mex += 1
        print(mex)
