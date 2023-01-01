# 负k次方根
# 求floor(a^(1/k))的值
# a<=2**64 k<=64


def kth_root_int(a: int, k: int) -> int:
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


T = int(input())
for _ in [0] * T:
    A, K = map(int, input().split())
    print(kth_root_int(A, K))
