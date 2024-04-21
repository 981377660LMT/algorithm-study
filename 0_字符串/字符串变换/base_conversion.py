# 进制转换(base conversion)

# base 進数から10進数に変える
# v[0]が最上位の桁を表す
# 6 = 110 なら v = {0, 1, 1}
from typing import List


def to_base_10(arr: List[int], fromBase: int) -> int:
    res = 0
    for x in arr:
        assert x >= 0
        assert x < fromBase
        res = res * fromBase + x
    return res


# 10 進数から toBase 進数に変える
def to_base_n(x: int, toBase: int) -> List[int]:
    if x == 0:
        return [0]
    res = []
    while x > 0:
        t = x % toBase
        res.append(t)
        x //= toBase
    return res[::-1]
