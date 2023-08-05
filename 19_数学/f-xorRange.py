# https://hitonanode.github.io/cplib-cpp/utilities/xor_of_interval.hpp
# 枚举异或范围(enumerateXorRange)
# lo <= (x ^ num) < hi となる x の区間 [L, R) の列挙

# 枚举所有的区间[L,R)使得在这个区间里的x满足
# !floor <= x^num < higher
# !注意每段[L,R)的长度都是2的幂次,这样的区间最多有O(log(higher-floor))个


from typing import Generator, Tuple


def xorRange(num: int, floor: int, higher: int) -> Generator[Tuple[int, int], None, None]:
    """
    ```
    floor <= (x ^ num) < higher
    ```
    """
    assert 0 <= floor, "floor must be non-negative"
    assert 0 <= num, "num must be non-negative"

    digit = 0
    while floor < higher:
        if floor & 1:
            yield (floor ^ num) << digit, ((floor ^ num) + 1) << digit
            floor += 1
        if higher & 1:
            higher -= 1
            yield (higher ^ num) << digit, ((higher ^ num) + 1) << digit
        floor >>= 1
        higher >>= 1
        num >>= 1
        digit += 1


if __name__ == "__main__":
    num, floor, higher = 3, 0, 10
    res1 = []
    for L, R in xorRange(num, floor, higher):
        res1.extend(range(L, R))
        print(L, R, "ok")
    res1.sort()

    res2 = []
    for i in range(1000):
        if floor <= (i ^ num) < higher:
            res2.append(i)

    assert res1 == res2
