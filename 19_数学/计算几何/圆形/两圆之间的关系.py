from typing import Tuple


Circle = Tuple[int, int, int]
Point = Tuple[int, int]

# !两圆之间的关系


def 外离(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) > (r1 + r2) * (r1 + r2)


def 外切(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) == (r1 + r2) * (r1 + r2)


def 相交(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) < (r1 + r2) * (r1 + r2)


def 内切(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) == (r1 - r2) * (r1 - r2)


def 内含(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) < (r1 - r2) * (r1 - r2)
