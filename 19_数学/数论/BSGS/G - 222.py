# 形如2,22,222,...的数列
# !这个数列第一个k的倍数的项是否存在, 若存在是第几项

# k<=1e8

# !等价于 2*(10^x-1)/9 ≡ 0 (mod k)
# !即 10^x ≡ 1 (mod k*9/gcd(k,2))
from math import gcd
from bsgs import exbsgs


# 即为扩展exbsgs
# TODO 有问题
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def find(k: int) -> int:
    return exbsgs(10, 1, k * 9 // gcd(k, 2))


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        k = int(input())
        print(find(k))
