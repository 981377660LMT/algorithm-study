# 0到x范围内二进制第k位为1的数的个数


from random import randint


def cal(upper: int, k: int) -> int:
    """[0, upper]中二进制第k(k>=0)位为1的数的个数.
    即满足 `num & (1 << k) > 0` 的数的个数
    """
    if k >= upper.bit_length():
        return 0
    res = upper // (1 << (k + 1)) * (1 << k)
    upper %= 1 << (k + 1)
    if upper >= 1 << k:
        res += upper - (1 << k) + 1
    return res


if __name__ == "__main__":

    def bruteForce(upper: int, k: int) -> int:
        res = 0
        for i in range(upper + 1):
            res += (i & (1 << k)) > 0
        return res

    for i in range(100):
        upper, k = randint(0, int(1e5)), randint(0, 20)
        if cal(upper, k) != bruteForce(upper, k):
            print(upper, k)
            print(cal(upper, k), bruteForce(upper, k))
            break
