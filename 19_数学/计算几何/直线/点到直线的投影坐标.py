from typing import Tuple


def projection(line: Tuple[int, int, int, int], point: Tuple[int, int]) -> Tuple[float, float]:
    """点在直线上的投影坐标."""
    x1, y1, x2, y2 = line
    x3, y3 = point
    dot = (x2 - x1) * (x3 - x1) + (y2 - y1) * (y3 - y1)
    d2 = ((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1)) ** 0.5  # hypot(x2 - x1, y2 - y1)
    bottom = dot / d2
    nx, ny = (x2 - x1) / d2, (y2 - y1) / d2
    return x1 + nx * bottom, y1 + ny * bottom


if __name__ == "__main__":
    x1, y1, x2, y2 = map(int, input().split())
    line = (x1, y1, x2, y2)
    q = int(input())
    for _ in range(q):
        x, y = map(int, input().split())
        point = (x, y)
        print(*projection(line, point))
