# https://atcoder.jp/contests/abc248/tasks/abc248_d
# RangeCountQuery


from collections import defaultdict
from typing import List, Tuple


def rangeCountQueryOffline(nums: List[int], queries: List[Tuple[int, int, int]]) -> List[int]:
    """
    O(n+q)离线区间频率查询.
    基于扫描线+前缀和实现.
    """
    n, q = len(nums), len(queries)
    res = [0] * q
    counter = defaultdict(lambda: 1)
    events = [[] for _ in range(n + 1)]
    for i, (start, end, target) in enumerate(queries):
        events[start].append((i, target))
        events[end].append((i, target))
    for i in range(n + 1):
        for j, target in events[i]:
            if not res[j]:
                res[j] -= counter[target]
            else:
                res[j] += counter[target]
        if i == n:
            break
        counter[nums[i]] += 1
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    nums = list(map(int, input().split()))
    q = int(input())
    queries = []
    for _ in range(q):
        start, end, target = map(int, input().split())
        start -= 1
        queries.append((start, end, target))
    res = rangeCountQueryOffline(nums, queries)
    print("\n".join(map(str, res)))
