def enumerateSubset(s: int):
    """降序枚举子集(包含s自身与空集)."""
    t = s
    while t:
        yield t
        t = (t - 1) & s
    yield 0


def enumerateSuperset(n: int, mask: int):
    """升序枚举超集

    Args:
        n: 集合大小
        mask: 初始二进制掩码
    """
    g1 = mask
    upper = 1 << n
    while g1 < upper:
        yield g1
        g1 = (g1 + 1) | mask


print(*enumerateSubset(0b101))  # [5,4,1,0]
print(*enumerateSuperset(4, 0b1001))  # [9,11,13,15]
