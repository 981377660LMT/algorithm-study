from collections import defaultdict
from typing import Generic, List, Optional, Tuple, TypeVar


INF = int(1e18)


V = TypeVar("V")


class EnumerateInterval(Generic[V]):
    """
    给定两个区间列表，每个区间列表都是成对 `不相交` 的。
    返回`[allStart, allEnd)`范围内的所有区间。
    返回值是一个列表，列表中的每个元素是一个五元组`[type, start, end, value1, value2]`。
    0: 不在两个区间列表中.
    1: 在第一个区间列表中,不在第二个区间列表中.
    2: 不在第一个区间列表中,在第二个区间列表中.
    3: 在两个区间列表中.
    """

    __slots__ = "_intervals1", "_intervals2"

    def __init__(self, intervals1: List[Tuple[int, int, V]], intervals2: List[Tuple[int, int, V]]):
        self._intervals1 = sorted(intervals1)
        self._intervals2 = sorted(intervals2)

    def getAll(self, start: int, end: int) -> List[Tuple[int, int, int, Optional[V], Optional[V]]]:
        n1, n2 = len(self._intervals1), len(self._intervals2)
        ptr1, ptr2 = 0, 0
        curStart = start
        res = []
        while ptr1 < n1 and self._intervals1[ptr1][1] <= curStart:
            ptr1 += 1
        while ptr2 < n2 and self._intervals2[ptr2][1] <= curStart:
            ptr2 += 1
        while curStart < end:
            start1 = self._intervals1[ptr1][0] if ptr1 < n1 else end
            end1 = self._intervals1[ptr1][1] if ptr1 < n1 else end
            start2 = self._intervals2[ptr2][0] if ptr2 < n2 else end
            end2 = self._intervals2[ptr2][1] if ptr2 < n2 else end
            intersect1 = start1 <= curStart < end1
            intersect2 = start2 <= curStart < end2
            if intersect1 and intersect2:
                minEnd = min2(end1, end2)
                res.append(
                    (3, curStart, minEnd, self._intervals1[ptr1][2], self._intervals2[ptr2][2])
                )
                curStart = minEnd
                if end1 == minEnd:
                    ptr1 += 1
                if end2 == minEnd:
                    ptr2 += 1
            elif intersect1:
                curEnd = min2(end1, start2)
                res.append((1, curStart, curEnd, self._intervals1[ptr1][2], None))
                curStart = curEnd
                if end1 == curEnd:
                    ptr1 += 1
            elif intersect2:
                curEnd = min2(end2, start1)
                res.append((2, curStart, curEnd, None, self._intervals2[ptr2][2]))
                curStart = curEnd
                if end2 == curEnd:
                    ptr2 += 1
            else:
                minStart = min2(start1, start2)
                res.append((0, curStart, minStart, None, None))
                curStart = minStart
        return res


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


if __name__ == "__main__":
    # 986. 区间列表的交集
    # https://leetcode.cn/problems/interval-list-intersections/
    class Solution:
        def intervalIntersection(
            self, firstList: List[List[int]], secondList: List[List[int]]
        ) -> List[List[int]]:
            intervals1 = [(start, end + 1, 1) for start, end in firstList]
            intervals2 = [(start, end + 1, 1) for start, end in secondList]
            E = EnumerateInterval(intervals1, intervals2)
            res = E.getAll(-INF, INF)
            return [[start, end - 1] for kind, start, end, *_ in res if kind == 3]

    # 两个游程编码区间的交集
    # https://atcoder.jp/contests/abc294/tasks/abc294_e
    # 给定两个长度均为L的数组a,b，求下标 i满足ai=bi的数量。
    # 每个游程编码数组以若干个二元组(v,l)依次表示，即有连续 l个值为 v，不断延续。
    # 例如，(1,2) 表示i=1和i=2的值为1
    def grid2xn(running1: List[Tuple[int, int]], running: List[Tuple[int, int]]) -> int:
        """两个游程编码区间的交集"""
        mp1 = defaultdict(list)  # v -> [(s1,e1,v1),(s2,e2,v1),...]
        mp2 = defaultdict(list)
        pre = 0
        for v, l in running1:
            cur = (pre, pre + l, v)
            mp1[v].append(cur)
            pre += l
        pre = 0
        for v, l in running2:
            cur = (pre, pre + l, v)
            mp2[v].append(cur)
            pre += l

        res = 0
        intersect = set(mp1.keys()) & set(mp2.keys())
        for v in intersect:
            intervals1 = mp1[v]
            intervals2 = mp2[v]
            E = EnumerateInterval(intervals1, intervals2)
            tmp = E.getAll(-INF, INF)
            res += sum(end - start for kind, start, end, *_ in tmp if kind == 3)

        return res

    L, n1, n2 = map(int, input().split())
    running1 = [tuple(map(int, input().split())) for _ in range(n1)]
    running2 = [tuple(map(int, input().split())) for _ in range(n2)]
    print(grid2xn(running1, running2))
