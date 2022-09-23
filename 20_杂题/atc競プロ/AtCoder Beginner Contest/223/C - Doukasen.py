"""导火线"""
# 有段导火线相连
# 第i个导火线的长度为Ai 燃烧速度为Bicm/s
# !两端同时点火, 求相遇时距离左端点的距离

# 每段保险丝的长度和燃烧速度都是不同的，
# 所以我们要用数学方法算出解析解是非常难的
# 二分答案

from typing import List, Tuple


EPS = 1e-7


def solve(n: int, fuses: List[Tuple[int, int]]):
    def calTime(nums: List[Tuple[int, int]], dist: float) -> float:
        res = 0.0
        for length, speed in nums:
            if dist <= length:
                res += dist / speed
                break
            else:
                res += length / speed
                dist -= length
        return res

    def check(mid: float) -> bool:
        """能否在mid处相遇"""
        leftTime, rightTime = calTime(fuses, mid), calTime(fuses[::-1], allSum - mid)
        return leftTime <= rightTime

    allSum = sum(fuse[0] for fuse in fuses)
    left, right = 0, allSum
    while left <= right:
        mid = (left + right) / 2
        if check(mid):
            left = mid + EPS
        else:
            right = mid - EPS
    return left


if __name__ == "__main__":
    n = int(input())
    fuses = [tuple(map(int, input().split())) for _ in range(n)]
    print(solve(n, fuses))
