# 区间等差数列常用操作
#
# arithmeticCount
# arithmeticSum
# findFloor
# findCeiling
# findFirst
# findLast


from typing import Tuple


def countOdds(lower: int, upper: int) -> int:
    """区间[lower, upper]内的奇数个数."""
    if lower > upper:
        return 0
    return arithmeticCount(lower, upper, 2, 1)


def countEvens(lower: int, upper: int) -> int:
    """区间[lower, upper]内的偶数个数."""
    if lower > upper:
        return 0
    return arithmeticCount(lower, upper, 2, 0)


def arithmeticCount(lower: int, upper: int, k: int, b: int) -> int:
    """区间[lower,upper]内形如k*x+b的个数."""
    if lower > upper:
        return 0
    if k == 0:
        return 1 if lower <= b <= upper else 0
    first, ok1 = findFirst(lower, upper, k, b)
    if not ok1:
        return 0
    last, ok2 = findLast(lower, upper, k, b)
    if not ok2:
        return 0
    return abs(last - first) // abs(k) + 1


def arithmeticSum(lower: int, upper: int, k: int, b: int) -> int:
    """区间[lower,upper]内形如k*x+b的和."""
    if lower > upper:
        return 0
    if k == 0:
        return b if lower <= b <= upper else 0
    first, ok1 = findFirst(lower, upper, k, b)
    if not ok1:
        return 0
    last, ok2 = findLast(lower, upper, k, b)
    if not ok2:
        return 0
    count = abs(last - first) // abs(k) + 1
    return (first + last) * count // 2


def findFloor(x: int, k: int, b: int) -> Tuple[int, bool]:
    """查找<=x的最大的形如k*x+b的数."""
    if k == 0:
        return (b, True) if b <= x else (0, False)
    step = abs(k)
    return (step * ((x - b) // step) + b, True)


def findCeiling(x: int, k: int, b: int) -> Tuple[int, bool]:
    """查找>=x的最小的形如k*x+b的数."""
    if k == 0:
        return (b, True) if b >= x else (0, False)
    step = abs(k)
    return (step * ((x - b + step - 1) // step) + b, True)


def findFirst(lower: int, upper: int, k: int, b: int) -> Tuple[int, bool]:
    """在区间[lower,upper]内查找第一个形如k*x+b的数."""
    if lower > upper:
        return (0, False)
    if k == 0:
        return (b, True) if lower <= b <= upper else (0, False)
    ceiling, ok1 = findCeiling(lower, k, b)
    if not ok1:
        return (0, False)
    if ceiling > upper:
        return (0, False)
    return (ceiling, True)


def findLast(lower: int, upper: int, k: int, b: int) -> Tuple[int, bool]:
    """在区间[lower,upper]内查找最后一个形如k*x+b的数."""
    if lower > upper:
        return (0, False)
    if k == 0:
        return (b, True) if lower <= b <= upper else (0, False)
    floor, ok1 = findFloor(upper, k, b)
    if not ok1:
        return (0, False)
    if floor < lower:
        return (0, False)
    return (floor, True)


if __name__ == "__main__":
    # 2894. 分类求和并作差
    # https://leetcode.cn/problems/divisible-and-non-divisible-sums-difference/
    class Solution:
        def differenceOfSums(self, n: int, m: int) -> int:
            sum_ = (1 + n) * n // 2
            okSum = arithmeticSum(1, n, m, 0)
            return sum_ - 2 * okSum

    def check() -> None:
        import random

        def findFloorBrute(x: int, k: int, b: int) -> Tuple[int, bool]:
            if k == 0:
                return (b, True) if b <= x else (0, False)
            cur = x
            while True:
                if (cur - b) % k == 0:
                    return cur, True
                cur -= 1

        def findCeilingBrute(x: int, k: int, b: int) -> Tuple[int, bool]:
            if k == 0:
                return (b, True) if b >= x else (0, False)
            cur = x
            while True:
                if (cur - b) % k == 0:
                    return cur, True
                cur += 1

        def findFirstBruteForce(lower: int, upper: int, k: int, b: int) -> Tuple[int, bool]:
            if lower > upper:
                return (0, False)
            if k == 0:
                return (b, True) if lower <= b <= upper else (0, False)
            for cur in range(lower, upper + 1):
                if (cur - b) % k == 0:
                    return cur, True
            return 0, False

        def findLastBruteForce(lower: int, upper: int, k: int, b: int) -> Tuple[int, bool]:
            if lower > upper:
                return (0, False)
            if k == 0:
                return (b, True) if lower <= b <= upper else (0, False)
            for cur in range(upper, lower - 1, -1):
                if (cur - b) % k == 0:
                    return cur, True
            return 0, False

        def arithmeticCountBruteForce(lower: int, upper: int, k: int, b: int) -> int:
            if k == 0:
                if b >= lower and b <= upper:
                    return 1
                return 0
            res = 0
            for cur in range(lower, upper + 1):
                if (cur - b) % k == 0:
                    res += 1
            return res

        def arithmeticSumBruteForce(lower: int, upper: int, k: int, b: int) -> int:
            if k == 0:
                if b >= lower and b <= upper:
                    return b
                return 0
            res = 0
            for cur in range(lower, upper + 1):
                if (cur - b) % k == 0:
                    res += cur
            return res

        def countOddsNaive(low: int, high: int) -> int:
            """区间[low, high]内的奇数个数."""
            return sum(1 for i in range(low, high + 1) if i & 1)

        def countEvensNaive(low: int, high: int) -> int:
            """区间[low, high]内的偶数个数."""
            return sum(1 for i in range(low, high + 1) if not i & 1)

        for _ in range(int(5e3)):
            x = random.randint(-int(5e3), int(5e3))
            k = random.randint(-int(5e3), int(5e3))
            b = random.randint(-int(5e3), int(5e3))
            res1, ok1 = findFloor(x, k, b)
            res2, ok2 = findFloorBrute(x, k, b)
            assert ok1 == ok2, f"{res1} {ok1} {res2} {ok2}"
            assert res1 == res2, f"{res1} {ok1} {res2} {ok2}"
            res1, ok1 = findCeiling(x, k, b)
            res2, ok2 = findCeilingBrute(x, k, b)
            assert ok1 == ok2, f"{res1} {ok1} {res2} {ok2}"
            assert res1 == res2, f"{res1} {ok1} {res2} {ok2}"
        for _ in range(int(5e3)):
            lower = random.randint(-int(5e3), int(5e3))
            upper = random.randint(-int(5e3), int(5e3))
            k = random.randint(-int(5e3), int(5e3))
            b = random.randint(-int(5e3), int(5e3))
            res1, ok1 = findFirst(lower, upper, k, b)
            res2, ok2 = findFirstBruteForce(lower, upper, k, b)
            assert ok1 == ok2, f"{res1} {ok1} {res2} {ok2}"
            assert res1 == res2, f"{res1} {ok1} {res2} {ok2}"
            res1, ok1 = findLast(lower, upper, k, b)
            res2, ok2 = findLastBruteForce(lower, upper, k, b)
            assert ok1 == ok2, f"{res1} {ok1} {res2} {ok2}"
            assert res1 == res2, f"{res1} {ok1} {res2} {ok2}"
        for _ in range(int(5e3)):
            lower = random.randint(-int(5e3), int(5e3))
            upper = random.randint(-int(5e3), int(5e3))
            k = random.randint(-int(5e3), int(5e3))
            b = random.randint(-int(5e3), int(5e3))
            res1 = arithmeticCount(lower, upper, k, b)
            res2 = arithmeticCountBruteForce(lower, upper, k, b)
            assert res1 == res2, f"{res1} {res2}"
            res1 = arithmeticSum(lower, upper, k, b)
            res2 = arithmeticSumBruteForce(lower, upper, k, b)
            assert res1 == res2, f"{res1} {res2}"
        for _ in range(int(5e3)):
            lower = random.randint(-int(5e3), int(5e3))
            upper = random.randint(-int(5e3), int(5e3))
            res1 = countOdds(lower, upper)
            res2 = countOddsNaive(lower, upper)
            assert res1 == res2, f"{res1} {res2}"
            res1 = countEvens(lower, upper)
            res2 = countEvensNaive(lower, upper)
            assert res1 == res2, f"{res1} {res2}"

        print("check success.")

    check()
