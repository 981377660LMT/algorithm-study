# https://atcoder.jp/contests/abc106/tasks/abc106_d

# 输入 n(≤500) m(≤2e5) q(≤1e5)。
# 然后输入 m 个在 [1,n] 内的闭区间，即每行输入两个数表示闭区间 [L,R]。
# 然后输入 q 个询问，每个询问也是输入两个数表示闭区间 [p,q]。
# 对每个询问，输出在 [p,q] 内的完整闭区间的个数。
# D - 每个查询包含完整区间的个数-容斥原理+区间dp/二维偏序


from collections import deque
import sys
from typing import List, Tuple

from sortedcontainers import SortedList

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def solve2(intervals: List[Tuple[int, int]], queries: List[Tuple[int, int]]) -> List[int]:
    """
    O(nlogn) 树状数组/有序集合(区间包含关系属于二维偏序)
    qlogq先对所有query排序,然后nlogn对于所有区间排序,把区间右侧端点放入树状数组(对应位置+1),
    然后从左往右遍历query,去掉所有左端点不满足的区间即可
    也算是离线查询
    """
    q = len(queries)
    intervals.sort()
    Q = sorted([(left, right, qi) for qi, (left, right) in enumerate(queries)])

    res = [0] * q
    ei = 0
    valid = deque()  # 有效区间
    sl = SortedList()  # 有效区间右端点
    for qLeft, qRight, qi in Q:
        # 1.加入所有合法的左端点
        while ei < len(intervals) and intervals[ei][0] <= qRight:
            sl.add(intervals[ei][1])
            valid.append(intervals[ei])
            ei += 1
        # 2. 删除不合法的左端点
        while valid and valid[0][0] < qLeft:
            sl.remove(valid[0][1])
            valid.popleft()
        res[qi] = sl.bisect_right(qRight)
    return res


if __name__ == "__main__":
    n, m, q = map(int, input().split())
    intervals = [tuple(map(int, input().split())) for _ in range(m)]
    queries = [tuple(map(int, input().split())) for _ in range(q)]
    # res = solve1(intervals, queries)
    res = solve2(intervals, queries)
    print(*res, sep="\n")
