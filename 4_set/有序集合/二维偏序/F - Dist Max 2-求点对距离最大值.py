"""
定义两点的距离为 `min(|x1-x2|, |y1-y2|)` 即魔改的切比雪夫距离
求点对距离最大值
n<=2e5


!最大化最小值:二分答案 
!二维偏序:一维排序+维护另一个维度前缀的最值

nlogn
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    points.sort()

    def check(mid: int) -> bool:
        """是否存在两个点使得 abs(x1-x2)>=mid 且 abs(y1-y2)>=mid

        !维护前缀的y坐标最大值和最小值,对每个点检查前缀中是否存在满足条件的点
        """
        minY, maxY = INF, -INF
        left = 0
        for right in range(n):
            while left <= right and points[right][0] - points[left][0] >= mid:
                minY = min(minY, points[left][1])
                maxY = max(maxY, points[left][1])
                left += 1
            if maxY - points[right][1] >= mid or points[right][1] - minY >= mid:
                return True
        return False

    left, right = 1, int(1e10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    print(right)
