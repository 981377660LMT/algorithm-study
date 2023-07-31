# 极角排序(Sort Points by Argument)
# n<=2e5 -1e9<=x,y<=1e9


from collections import defaultdict, deque
from typing import List, Tuple


def angleSort(points: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    """极角排序，返回排序后的结果.

    按照atan2(x,y) 排序
    从-pi开始,逆时针到pi
    eg: (-10,-1) -> (0, -1) -> (1, 0) -> (0, 1) -> (-1, 0)

    atan2(x<0,y=0) = pi
    atan2(0,0) = 0
    """

    def msort(xx):
        if not xx:
            return
        n = len(xx)
        a = [xx, xx[:]]
        lr = [(0, n, 1, 0)]
        while lr:
            l, r, f, g = lr.pop()
            m = (l + r) // 2
            if f:
                lr.append((l, r, 0, g))
                if m - l > 1:
                    lr.append((l, m, 1, g ^ 1))
                if r - m > 1:
                    lr.append((m, r, 1, g ^ 1))
            else:
                i, j, p, q = l, m, m - 1, r - 1
                a1 = a[g]
                a2 = a[g ^ 1]
                for k in range((r - l) // 2):
                    x, y = a2[i]
                    s, t = a2[j]
                    if s * y - t * x > 0:
                        a1[l + k] = a2[j]
                        j += 1
                    else:
                        a1[l + k] = a2[i]
                        i += 1
                    x, y = a2[p]
                    s, t = a2[q]
                    if s * y - t * x > 0:
                        a1[r - 1 - k] = a2[p]
                        p -= 1
                    else:
                        a1[r - 1 - k] = a2[q]
                        q -= 1
                if (r - l) & 1:
                    a1[m] = a2[i] if i == p else a2[j]

    z1, z2, z3, z4, z5 = [], [], [], [], []
    for x, y in points:
        if x == y == 0:
            z5.append((x, y))
        elif y >= 0:
            if x >= 0:
                z1.append((x, y))
            else:
                z2.append((x, y))
        else:
            if x < 0:
                z3.append((x, y))
            else:
                z4.append((x, y))

    msort(z1)
    msort(z2)
    msort(z3)
    msort(z4)

    ids = z3 + z4 + z5 + z1 + z2
    return ids


def angleArgSort(points: List[Tuple[int, int]]) -> List[int]:
    """极角排序，返回每个点在极角排序后的索引."""
    res = angleSort(points)
    mp = defaultdict(list)
    for i, p in enumerate(res):
        mp[p].append(i)
    order = [0] * len(points)
    for i, p in enumerate(points):
        order[i] = mp[p].pop()
    return order  # type: ignore


if __name__ == "__main__":
    res = angleSort(
        [
            (1, 0),
            (0, 0),
            (-1, 0),
            (0, 1),
            (0, -1),
            (1, 1),
            (2, 2),
            (-10, -1),
        ]
    )
    print(res)

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    res = angleSort(points)
    for x, y in res:
        print(x, y)
# 8
# 1 0
# 0 0
# -1 0
# 0 1
# 0 -1
# 1 1
# 2 2
# -10 -1

# 输出
# -10 -1
# 0 -1
# 1 0
# 0 0
# 1 1
# 2 2
# 0 1
# -1 0
