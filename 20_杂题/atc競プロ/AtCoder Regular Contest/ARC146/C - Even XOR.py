from functools import reduce
from operator import xor
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 0-2^n-1
# 奇数个 或者 异或不为0
# 容斥?
# 奇数个元素
# 异或为0的集合个数


# n = int(input())
# res1 = pow(2, n - 1, MOD)
# res2 = 0

from itertools import chain, combinations
from typing import Collection, Generator, List, TypeVar


# 2. powerset 顺序枚举
T = TypeVar("T")


def powerset(collection: Collection[T], isAll=True):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) + 1 if isAll else len(collection)
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


# 奇数个元素异或怎么为0
for n in range(1, 5):
    res = 0
    # subset
    for subset in powerset(range(2**n)):
        if not subset:
            res += 1
            continue
        if len(subset) & 1 or reduce(xor, subset) != 0:
            res += 1
    print(res)
