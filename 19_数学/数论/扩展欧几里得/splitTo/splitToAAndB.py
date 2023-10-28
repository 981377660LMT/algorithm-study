from typing import Tuple


def splitToAAndB(num: int, a: int, b: int, minimize=True) -> Tuple[int, int, bool]:
    """
    将 num 拆分成 a 和 b 的和，使得拆分的个数最(多/少).

    :param num: 正整数.
    :param a: 正整数.
    :param b: 正整数.
    :param minimize: 是否使得拆分的个数最少. 默认为最少(True).

    :returns: [countA, countB, ok] countA和countB分别是拆分成a和b的个数,ok表示是否可以拆分.
    """
    n, x1, y1, x2, y2 = solveLinearEquation(a, b, num, allowZero=True)
    if n < 0:
        return 0, 0, False
    if n > 0:
        res1Smaller = x1 + y1 < x2 + y2
        return (x1, y1, True) if res1Smaller == minimize else (x2, y2, True)

    # 存在整数解但不存在正整数解，检查其中一项是否可以为0
    modA, modB = num % a, num % b
    if modA and modB:
        return 0, 0, False
    if modA:
        return 0, num // b, True
    if modB:
        return num // a, 0, True
    div1, div2 = num // a, num // b
    res1Smaller = div1 < div2
    return (div1, 0, True) if res1Smaller == minimize else (0, div2, True)


def solveLinearEquation(a: int, b: int, c: int, allowZero=False) -> Tuple[int, int, int, int, int]:
    """
    a*x + b*y = c 的通解为
    x = (c/g)*x0 + (b/g)*k
    y = (c/g)*y0 - (a/g)*k
    其中 g = gcd(a,b) 且需要满足 g|c
    x0 和 y0 是 ax+by=g 的一组特解（即 exgcd(a,b) 的返回值）

    为方便讨论，这里要求输入的 a b c 必须为正整数

    返回: 正整数解的个数（无解时为 -1,无正整数解时为 0)
          x 取最小正整数时的解 x1 y1,此时 y1 是最大正整数解
          y 取最小正整数时的解 x2 y2,此时 x2 是最大正整数解
    """
    g, x0, y0 = exgcd(a, b)

    # 无解
    if c % g != 0:
        return -1, 0, 0, 0, 0

    a //= g
    b //= g
    c //= g
    x0 *= c
    y0 *= c

    x1 = x0 % b
    if allowZero:
        if x1 < 0:
            x1 += b
    else:
        if x1 <= 0:
            x1 += b
    k1 = (x1 - x0) // b
    y1 = y0 - k1 * a

    y2 = y0 % a
    if allowZero:
        if y2 < 0:
            y2 += a
    else:
        if y2 <= 0:
            y2 += a
    k2 = (y0 - y2) // a
    x2 = x0 + k2 * b

    # 无正整数解
    if y1 <= 0:
        return 0, x1, y1, x2, y2

    # k 越大 x 越大
    return k2 - k1 + 1, x1, y1, x2, y2


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """求解二元一次不定方程 ax+by=gcd(a,b) 的特解(x,y),特解满足 |x|<=|b|, |y|<=|a|."""
    if b == 0:
        return a, 1, 0
    gcd_, y, x = exgcd(b, a % b)
    y -= a // b * x
    return gcd_, x, y


if __name__ == "__main__":
    from typing import List
    from collections import Counter

    class Solution:
        def minGroupsForValidAssignment(self, nums: List[int]) -> int:
            freqCounter = Counter(Counter(nums).values())
            res = len(nums)
            for size in range(1, len(nums) + 1):
                tmp = 0
                for num in freqCounter:
                    count1, count2, ok = splitToAAndB(num, size, size + 1, minimize=True)
                    if not ok:
                        break
                    tmp += (count1 + count2) * freqCounter[num]
                else:
                    res = min(res, tmp)
            return res
