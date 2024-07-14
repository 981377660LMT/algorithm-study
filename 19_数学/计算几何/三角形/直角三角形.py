def isRightTriangle(x1: int, y1: int, x2: int, y2: int, x3: int, y3: int) -> bool:
    """判断是否为直角三角形."""

    def f(x1: int, y1: int, x2: int, y2: int, x3: int, y3: int) -> bool:
        ij = (x1 - x2, y1 - y2)
        ik = (x3 - x2, y3 - y2)
        return ij[0] * ik[0] + ij[1] * ik[1] == 0

    return f(x1, y1, x2, y2, x3, y3) or f(x2, y2, x3, y3, x1, y1) or f(x3, y3, x1, y1, x2, y2)


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc362/tasks/abc362_b
    x1, y1 = map(int, input().split())
    x2, y2 = map(int, input().split())
    x3, y3 = map(int, input().split())
    if isRightTriangle(x1, y1, x2, y2, x3, y3):
        print("Yes")
    else:
        print("No")
