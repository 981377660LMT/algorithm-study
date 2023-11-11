from functools import cmp_to_key
from typing import List, Tuple


def angleArgSort(points: List[Tuple[int, int]]) -> List[int]:
    """极角排序，返回值为点的下标."""

    def compare(i: int, j: int) -> int:
        x1, y1 = points[i]
        x2, y2 = points[j]
        return x2 * y1 - x1 * y2

    lower, origin, upper = [], [], []
    zero = (0, 0)
    for i, p in enumerate(points):
        if p == zero:
            origin.append(i)
        elif p[1] < 0 or (p[1] == 0 and p[0] > 0):
            lower.append(i)
        else:
            upper.append(i)

    lower.sort(key=cmp_to_key(compare))
    upper.sort(key=cmp_to_key(compare))
    return lower + origin + upper


def angleSort(points: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    """极角排序，返回值为点的坐标."""
    indices = angleArgSort(points)
    return [points[i] for i in indices]


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    res = angleSort(points)
    for x, y in res:
        print(x, y)
