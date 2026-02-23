# https://yukicoder.me/problems/no/1339
# 求有理数1/n的循环节长度

# !num % n => num * 10  => num % n => num * 10 => ...
# 离散对数 求 10**x ≡ 1 (mod n) 的最小解x

from math import ceil, sqrt


def bsgs(base: int, target: int, p: int) -> int:
    """Baby-step Giant-step

    在base和p互质的情况下,求解 base^x ≡ target (mod p) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(p)))

    https://dianhsu.com/2022/08/27/template-math/#bsgs
    """
    target %= p
    mp = dict()
    t = ceil(sqrt(p))
    val = 1
    for i in range(t):
        tv = target * val % p
        mp[tv] = i
        val = val * base % p

    base, val = val, 1
    if base == 0:
        return 1 if target == 0 else -1

    for i in range(t + 1):
        tv = mp.get(val, -1)
        if tv != -1 and i * t - tv > 0:  # !注意这里取等号表示允许最小解为0
            return i * t - tv
        val = val * base % p

    return -1


def solve(n: int) -> int:
    while n % 2 == 0:
        n //= 2
    while n % 5 == 0:
        n //= 5
    return bsgs(10, 1, n)


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        print(solve(int(input())))
