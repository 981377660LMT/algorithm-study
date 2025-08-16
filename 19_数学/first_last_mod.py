def first_last_mod(start: int, end: int, mod: int, remainder: int):
    """
    在[start, end)区间内，寻找第一个和最后一个满足x % mod == remainder的数。
    如果不存在，返回None, None。
    """
    if start >= end:
        return None, None
    if not (0 <= remainder < mod):
        return None, None
    r = start % mod
    delta = (remainder - r) % mod
    first = start + delta
    if first >= end:
        return None, None
    last = first + ((end - 1 - first) // mod) * mod
    return first, last


if __name__ == "__main__":
    assert first_last_mod(1, 10, 3, 2) == (2, 8)
    assert first_last_mod(1, 10, 3, 1) == (1, 7)
    assert first_last_mod(1, 11, 3, 1) == (1, 10)
    assert first_last_mod(1, 10, 3, 0) == (3, 9)
    # 边界测试
    assert first_last_mod(0, 0, 3, 0) == (None, None)
    assert first_last_mod(5, 5, 3, 2) == (None, None)
    assert first_last_mod(10, 20, 1, 0) == (10, 19)
    assert first_last_mod(7, 8, 2, 1) == (7, 7)
    assert first_last_mod(7, 8, 2, 0) == (None, None)
    # 非法参数
    assert first_last_mod(1, 10, 3, 3) == (None, None)
    assert first_last_mod(1, 10, 3, -1) == (None, None)
    print("All tests passed.")
