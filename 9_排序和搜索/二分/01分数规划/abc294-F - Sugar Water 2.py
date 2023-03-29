# https://atcoder.jp/contests/abc294/tasks/abc294_f

# 两排糖水，糖成份x和水成分 y。
# !从两排中各取一个糖水混合，得糖浓度 (x1+x2)/(x1+y1+x2+y2)。
# 问所有方案中，糖浓度第 k大是多少(1<=k<=m*n)。
# m,n<=5e4

# !01分数规划=> 二分答案+公式变形
# (xi+xj)<=mid*(xi+yi+xj+yj)
# !将i和j分离得到
# (mid-1)*xi +mid*yi + (mid-1)*xj +mid*yj >=0
# !把(mid-1)*x+mid*y作为新的数组,也就是求两个数组中和大于等于0的对数 (排序+双指针或者二分都以)

from typing import List, Tuple


EPS = 1e-12


def sugarWater2(pairs1: List[Tuple[int, int]], pairs2: List[Tuple[int, int]], kthMin: int) -> float:
    def countNGT(mid: float) -> int:
        """有多少个不超过mid的候选"""
        A = [a * (mid - 1) + b * mid for a, b in pairs1]
        B = [a * (mid - 1) + b * mid for a, b in pairs2]
        A.sort()
        B.sort()
        res, right = 0, len(B) - 1
        for v in A:
            while right >= 0 and v + B[right] >= 0:
                right -= 1
            res += len(B) - right - 1
        return res

    left, right = 0, 1
    while left <= right:
        mid = (left + right) / 2
        if countNGT(mid) < kthMin:
            left = mid + EPS
        else:
            right = mid - EPS
    return left


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    pairs1 = [tuple(map(int, input().split())) for _ in range(n)]
    pairs2 = [tuple(map(int, input().split())) for _ in range(m)]
    kthMin = n * m + 1 - k
    print(sugarWater2(pairs1, pairs2, kthMin) * 100)
