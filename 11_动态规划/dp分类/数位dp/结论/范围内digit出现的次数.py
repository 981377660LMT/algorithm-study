# https://www.acwing.com/problem/content/340/
#  统计区间[L,R]出现0123456789的各个数字总次数
# 每个结果包含十个用空格隔开的数字，
# 第一个数字表示 0 出现的次数，第二个数字表示 1 出现的次数，以此类推。
# 1≤a,b≤231−1,
# 1≤N<100
# https://www.acwing.com/activity/content/code/content/4041182/

# 范围内的数字中digit出现的次数


def countDigit(n: int, digit: int, base=10) -> int:
    """O(lgn)求[1,n]中digit出现的次数."""
    res = 0
    left, right = 0, 0

    def calLen(n: int) -> int:
        res = 0
        while n:
            res += 1
            n //= base
        return res

    len_ = calLen(n)
    for i in range(1, len_ + 1):
        right = base ** (i - 1)
        left = n // (right * base)
        res += left * right if digit else (left - 1) * right
        d = (n // right) % base
        if d == digit:
            res += n % right + 1
        elif d > digit:
            res += right
    return res


def countDigitIn(left: int, right: int, digit: int, base=10) -> int:
    """统计[left,right]中digit出现的次数."""
    if left > right:
        return 0
    if left == 0:
        return countDigit(right, digit, base=base) + (digit == 0)
    return countDigit(right, digit, base=base) - countDigit(left - 1, digit, base=base)


if __name__ == "__main__":
    print(countDigitIn(1, 12, 1, base=2))
    print(sum(v.bit_count() for v in range(1, 13)))
