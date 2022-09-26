"""
苹果围成一圈 从1开始 每次可以在每个位置吃一个苹果
求吃完k个苹果时 每个位置剩下的苹果数

二分圈数 O(nlogk)
!找到最大的right使得恰好不能吃完k个苹果
剩下的再走一圈均摊
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k = map(int, input().split())
    apples = list(map(int, input().split()))

    def check(mid: int) -> bool:
        """mid圈不能吃完k个苹果"""
        res = 0
        for i in range(n):
            res += min(mid, apples[i])
        return res < k

    left, right = 0, int(1e12) + 5
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    # !找到最大的right使得恰好不能吃完k个苹果
    remain = k
    nums = apples[:]
    for i in range(n):
        eat = min(right, nums[i])
        nums[i] -= eat
        remain -= eat

    for i in range(n):
        if remain == 0:
            break
        if nums[i]:
            nums[i] -= 1
            remain -= 1

    print(*nums)
