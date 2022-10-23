"""
Range Affine Range Sum
区间映射/区间求和/奇妙序列

q个操作:
0 left right mul add: [left:right]切片内的每个数都乘以b, 再加上c
1 left right: [left:right]切片内的所有数的和 mod 998244353
"""

import sys
from AtcoderLazySegmentTree import AtcoderLazySegmentTree


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# !64位
S = int  # !sum用前32位, length用后32位 (sum,length)
F = int  # !mul用前32位, add用后32位  (mul,add)
POW32 = 1 << 32


def e() -> S:
    return 0


def id() -> "F":
    return POW32


def op(leftData: "S", rightData: "S") -> "S":
    leftSum, leftLen = leftData // POW32, leftData % POW32
    rightSum, rightLen = rightData // POW32, rightData % POW32
    return ((leftSum + rightSum) % MOD) * POW32 + (leftLen + rightLen)


def mapping(parentLazy: "F", childData: "S") -> "S":
    lazyMul, lazyAdd = parentLazy // POW32, parentLazy % POW32
    childSum, childLen = childData // POW32, childData % POW32
    return ((lazyMul * childSum + lazyAdd * childLen) % MOD) * POW32 + childLen


def composition(parentLazy: "F", childLazy: "F") -> "F":
    lazyMul, lazyAdd = parentLazy // POW32, parentLazy % POW32
    childSum, childLen = childLazy // POW32, childLazy % POW32
    return ((lazyMul * childSum) % MOD) * POW32 + ((lazyMul * childLen + lazyAdd) % MOD)


if __name__ == "__main__":
    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    init = [(num << 32) + 1 for num in nums]
    tree = AtcoderLazySegmentTree(init, e=e, id=id, op=op, mapping=mapping, composition=composition)
    for _ in range(q):
        kind, *args = map(int, input().split())
        if kind == 0:
            left, right, mul, add = args
            tree.update(left, right, (mul << 32) + add)
        else:
            left, right = args
            print((tree.query(left, right) >> 32) % MOD)
