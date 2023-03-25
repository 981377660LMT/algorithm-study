# 400. 第 N 位数字
# https://leetcode.cn/problems/nth-digit/
# 给你一个整数 n ，请你在无限的整数序列 [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ...] 中找出并返回第 n 位上的数字。


def findNthDigit(n: int) -> int:
    """O(lgn)返回第n位数码, n从1开始."""
    base, count = 1, 1
    while True:
        digitSum = count * 9 * base
        if n > digitSum:
            n -= digitSum
            base *= 10
            count += 1
        else:
            break
    n -= 1  # 从0开始
    num, mod = base + n // count, n % count
    return int(str(num)[mod])


def countDigit(n: int, digit: int) -> int:
    """O(lgn)求[1,n]中digit出现的次数."""
    res = 0
    left, right = 0, 0
    len_ = len(str(n))
    for i in range(1, len_ + 1):
        right = 10 ** (i - 1)
        left = n // (right * 10)
        res += left * right if digit else (left - 1) * right
        d = (n // right) % 10
        if d == digit:
            res += n % right + 1
        elif d > digit:
            res += right
    return res


def countDigitIn(left: int, right: int, digit: int) -> int:
    """统计[left,right]中digit出现的次数."""
    if left > right:
        return 0
    if left == 0:
        return countDigit(right, digit) + (digit == 0)
    return countDigit(right, digit) - countDigit(left - 1, digit)


def digitsLen(x: int) -> int:
    """所有x位数的位数之和."""
    return x * 9 * 10 ** (x - 1)


if __name__ == "__main__":
    # countDigit
    for left in range(200):
        for right in range(left, 200):
            for j in range(10):
                sum1 = countDigitIn(left, right, j)
                sum2 = sum(str(i).count(str(j)) for i in range(left, right + 1))
                assert sum1 == sum2, f"{left}, {right}, {j}, {sum1}, {sum2}"
    print("countDigit passed")
