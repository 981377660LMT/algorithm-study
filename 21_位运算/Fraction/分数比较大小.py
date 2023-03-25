# 斜率比较大小 (a1/b1) 和 (a2/b2) 谁大, 分母不为0
# 分数比较大小 compareFraction


def eq(a1: int, b1: int, a2: int, b2: int) -> bool:
    """a1/b1 == a2/b2"""
    return a1 * b2 == a2 * b1


def lt(a1: int, b1: int, a2: int, b2: int) -> bool:
    """a1/b1 < a2/b2"""
    diff = a1 * b2 - a2 * b1
    mul = b1 * b2
    if diff == 0:
        return False
    return (diff > 0) ^ (mul > 0)


def le(a1: int, b1: int, a2: int, b2: int) -> bool:
    """a1/b1 <= a2/b2"""
    diff = a1 * b2 - a2 * b1
    mul = b1 * b2
    if diff == 0:
        return True
    return (diff > 0) ^ (mul > 0)


if __name__ == "__main__":
    from random import randint

    for _ in range(100000):
        a1, b1, c1, d1 = (
            randint(-100, 100),
            randint(-100, 100),
            randint(-100, 100),
            randint(-100, 100),
        )
        if b1 == 0 or d1 == 0:
            continue
        assert lt(a1, b1, c1, d1) == (a1 / b1 < c1 / d1) == (not le(c1, d1, a1, b1))
    print("ok")
