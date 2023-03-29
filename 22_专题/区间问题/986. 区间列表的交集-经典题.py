# 两个区间列表的交集(区间交集)

from collections import defaultdict
from typing import List, Tuple


Interval = Tuple[int, int]


def solve(intervals1: List[Interval], intervals2: List[Interval]) -> List[Interval]:
    """双指针求两个已排序的区间列表的交集"""
    n1, n2 = len(intervals1), len(intervals2)
    res = []
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1, s2, e2 = *intervals1[left], *intervals2[right]
        # 相交
        if s1 <= e2 <= e1 or s2 <= e1 <= e2:
            # 尽量往内缩
            res.append((max(s1, s2), min(e1, e2)))
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res


def solve2(intervals1: List[Interval], intervals2: List[Interval]) -> int:
    """两个已排序的区间列表相交的长度之和"""
    n1, n2 = len(intervals1), len(intervals2)
    res = 0
    left, right = 0, 0
    while left < n1 and right < n2:
        s1, e1, s2, e2 = *intervals1[left], *intervals2[right]
        if s1 <= e2 <= e1 or s2 <= e1 <= e2:  # 相交
            res += min(e1, e2) - max(s1, s2) + 1  # [1,1] 区间长度为1
        if e1 < e2:
            left += 1
        else:
            right += 1
    return res


# 986. 两个游程编码区间的交集
# https://atcoder.jp/contests/abc294/tasks/abc294_e
# 给定两个长度均为L的数组a,b，求下标 i满足ai=bi的数量。
# 每个游程编码数组以若干个二元组(v,l)依次表示，即有连续 l个值为 v，不断延续。
# 例如，(1,2) 表示i=1和i=2的值为1
def grid2xn(running1: List[Tuple[int, int]], running: List[Tuple[int, int]]) -> int:
    """两个游程编码区间的交集"""
    mp1 = defaultdict(list)  # v -> [(s1,e1),(s2,e2),...]
    mp2 = defaultdict(list)
    pre = 0
    for v, l in running1:
        cur = (pre, pre + l - 1)
        mp1[v].append(cur)
        pre += l
    pre = 0
    for v, l in running2:
        cur = (pre, pre + l - 1)
        mp2[v].append(cur)
        pre += l

    res = 0
    intersect = set(mp1.keys()) & set(mp2.keys())
    for v in intersect:
        res += solve2(mp1[v], mp2[v])
    return res


if __name__ == "__main__":
    L, n1, n2 = map(int, input().split())
    running1 = [tuple(map(int, input().split())) for _ in range(n1)]
    running2 = [tuple(map(int, input().split())) for _ in range(n2)]
    print(grid2xn(running1, running2))
