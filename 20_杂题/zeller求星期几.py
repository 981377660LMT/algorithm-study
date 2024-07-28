# zeller求星期几
# 给定年月日,求星期数


def zeller(y: int, m: int, d: int) -> int:
    if m <= 2:
        m += 12
        y -= 1
    return (y + y // 4 - y // 100 + y // 400 + (13 * m + 8) // 5 + d) % 7


assert (zeller(2023, 2, 4)) == 6
