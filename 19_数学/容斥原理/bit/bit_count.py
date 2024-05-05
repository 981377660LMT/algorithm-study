#  bit_count


import sys


LONG_LONG_UPPER = 1 << 63 - 1


def bit_count(n: int) -> int:
    if sys.version_info >= (3, 10):
        return n.bit_count()  # 非常快
    if n <= LONG_LONG_UPPER:
        return long_long_bit_count(n)  # 快
    return int_bit_count(n)  # 一般快


def long_long_bit_count(n: int) -> int:
    """64位数字的bit_count  334 ms"""
    c = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
    c = (c & 0x3333333333333333) + ((c >> 2) & 0x3333333333333333)
    c = (c & 0x0F0F0F0F0F0F0F0F) + ((c >> 4) & 0x0F0F0F0F0F0F0F0F)
    c = (c & 0x00FF00FF00FF00FF) + ((c >> 8) & 0x00FF00FF00FF00FF)
    c = (c & 0x0000FFFF0000FFFF) + ((c >> 16) & 0x0000FFFF0000FFFF)
    c = (c & 0x00000000FFFFFFFF) + ((c >> 32) & 0x00000000FFFFFFFF)
    return c


def int_bit_count(n: int) -> int:
    """471 ms"""
    res = 0
    while n:
        n &= n - 1
        res += 1
    return res


def bin_str_bit_count(n: int) -> int:
    """1461 ms"""
    return bin(n).count("1")


num = (1 << 64) - 1  # 最多64位
print(long_long_bit_count(num))
print(int_bit_count(num))
print(bin_str_bit_count(num))
