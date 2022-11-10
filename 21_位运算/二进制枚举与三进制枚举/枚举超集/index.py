from typing import List


def enumerateSubset(mask: int) -> List[int]:
    """枚举(非空)子集"""
    res = []
    g1 = mask
    while g1 > 0:
        res.append(g1)
        g1 = (g1 - 1) & mask
    return res


def enumerateSuperset(n: int, mask: int) -> List[int]:
    """枚举超集

    Args:
        n: 集合大小
        mask: 初始二进制掩码
    """
    res = []
    g1 = mask
    upper = 1 << n
    while g1 < upper:
        res.append(g1)
        g1 = (g1 + 1) | mask
    return res


print(enumerateSubset(0b101))  # [5,4,1]
print(enumerateSuperset(4, 0b1001))  # [9,11,13,15]
