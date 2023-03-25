# 负k次方根(kthRoot/floorRoot)/单位根
# 求floor(a^(1/k))的值
# a<=2**64 k<=64


from typing import Optional


def floorRoot(a: int, k: int) -> int:
    """floor(a^(1/k))"""
    assert 0 <= a and 0 < k
    if a == 0:
        return 0
    if k == 1:
        return a
    res = int(pow(a, 1 / k))
    while pow(res + 1, k) <= a:
        res += 1
    while pow(res, k) > a:
        res -= 1
    return res


def floorRoot2(a: int, k: int) -> int:
    if a <= 0:
        return 0
    r1 = k - 1
    x = int(a ** (1.0 / k) * (1 + 1e-12))
    while True:
        y = (r1 * x + a // (x**r1)) // k
        if y >= x:
            return x
        x = y


def ceilRoot(a: int, k: int) -> int:
    """ceil(a^(1/k))"""
    assert 0 <= a and 0 < k
    if a == 0:
        return 0
    if k == 1:
        return a
    res = int(pow(a, 1 / k)) + 1
    while pow(res, k) < a:
        res += 1
    while a <= pow(res - 1, k):
        res -= 1
    return res


def kthPower(a: int, k: int) -> Optional[int]:
    """判断a是否为一个整数的k次方
    如果是,则返回满足b^k=a的b,否则返回None.
    """
    abs_ = abs(a)
    sign = a // abs_ if a else 0
    b = floorRoot(abs_, k)
    if pow(sign * b, k) == a:
        return sign * b


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/1666
    # https://yukicoder.me/submissions/695435
    # 求出第k个形如a^b的数(累乘数) ,a>=1,b>=2
    # 1<=k<=1e9
    def solve(k: int) -> int:
        ...
