# 统计区间内的个数


def divCount(lower: int, upper: int, k: int) -> int:
    """区间[lower,upper]内k的倍数个数."""
    return upper // k - (lower - 1) // k


def divSum(lower: int, upper: int, k: int) -> int:
    """区间[lower,upper]内k的倍数和."""
    if lower > upper:
        return 0
    assert k > 0, "k must be positive"

    def f(right: int) -> int:
        if right < k:
            return 0
        first = k
        last = right // k * k
        count = (last - first) // k + 1
        return (first + last) * count // 2

    return f(upper) - f(lower - 1)


if __name__ == "__main__":
    print(divCount(1, 10, 3))
    print(divSum(1, 10, 3))

    # 2894. 分类求和并作差
    # https://leetcode.cn/problems/divisible-and-non-divisible-sums-difference/
    class Solution:
        def differenceOfSums(self, n: int, m: int) -> int:
            sum_ = (1 + n) * n // 2
            okSum = divSum(1, n, m)
            return sum_ - 2 * okSum
