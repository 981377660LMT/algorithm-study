# https://maspypy.github.io/library/enumerate/xor_range.hpp
# 枚举异或范围
# lo <= (x ^ a) < hi となる x の区間 [L, R) の列挙

# 枚举所有的区间[L,R)使得在这个区间里的x满足
# lower<=x^a<upper


def xorRange(a: int, lower: int, upper: int):
    """
    ```
    lower <= (x ^ a) < upper
    ```
    """
    for k in range(64):
        if lower == upper:
            break
        b = 1 << k
        if lower & b:
            yield lower ^ a, (lower ^ a) + b
            lower += b
        if upper & b:
            yield (upper - b) ^ a, ((upper - b) ^ a) + b
            upper -= b
        if a & b:
            a ^= b


if __name__ == "__main__":
    for L, R in xorRange(3, 0, 10):
        print(L, R)
    res = []
    for i in range(100):
        if 0 <= i ^ 3 < 10:
            res.append(i)
    print(res)
