from typing import Tuple


Circle = Tuple[int, int, int]


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


if __name__ == "__main__":
    circle1 = tuple(map(int, input().split()))  # x, y, r
    circle2 = tuple(map(int, input().split()))  # x, y, r
    if 内含(circle1, circle2):
        print(1)
    elif 内切(circle1, circle2):
        print(2)
    elif 相交(circle1, circle2):
        print(3)
    elif 外切(circle1, circle2):
        print(4)
    elif 外离(circle1, circle2):
        print(5)
