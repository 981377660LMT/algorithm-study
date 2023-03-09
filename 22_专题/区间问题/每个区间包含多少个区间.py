# https://www.luogu.com.cn/problem/CF652D

# 给定n个区间，求每个区间包含多少个其他的区间。
# !二维偏序 把区间(l,r)想象成二维平面上的一个点
# n<=2e5
# 排序

from typing import List, Tuple
from bisect import bisect_left

from sortedcontainers import SortedList

INF = int(1e18)


def nestedSegments(segments: List[Tuple[int, int]]) -> List[int]:
    qs = [(x, y, i) for i, (x, y) in enumerate(segments)]
    qs.sort(key=lambda x: (x[1], -x[0]))
    sl = SortedList()  # 维护x的有序集合
    res = [0] * len(segments)
    ei = 0
    while ei < len(qs):
        group = [qs[ei][2]]
        x, y = qs[ei][0], qs[ei][1]
        while ei + 1 < len(qs) and qs[ei + 1][0] == x and qs[ei + 1][1] == y:
            ei += 1
            group.append(qs[ei][2])
        for _ in range(len(group)):
            sl.add(x)
        more = len(sl) - bisect_left(sl, x)  # >=x的个数
        for i in group:
            res[i] = more - 1
        ei += 1

    return res


if __name__ == "__main__":
    n = int(input())
    segments = list(tuple(map(int, input().split())) for _ in range(n))
    res = nestedSegments(segments)
    print(*res, sep="\n")
