# E - Maximize XOR
# https://atcoder.jp/contests/abc386/tasks/abc386_e
#
# 组合数性质
# 给你一个长度为n的数组，再给你一个数字k，保证从n个数字里面选择k个数的总数小于1e6,
# 现在让你从数组n中选择k个数字，使得这k个数字的异或值最大，求出这个最大的异或值.
# C(n,k) <= 1e6, n,k<=2e.
#
# 爆搜+剪枝：当剩下的数字不足当前的k值时候，直接把剩下的全部选上就行.

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def maximizeXor(nums: List[int], k: int) -> int:
    n = len(nums)
    sufXor = [0] * (n + 1)
    for i in range(n - 1, -1, -1):
        sufXor[i] = sufXor[i + 1] ^ nums[i]

    res = 0

    def bt(index: int, remain: int, xor: int) -> None:
        nonlocal res
        if remain == 0:
            res = max2(res, xor)
            return
        if index + remain >= n:
            res = max2(res, xor ^ sufXor[index])
            return
        bt(index + 1, remain, xor)
        bt(index + 1, remain - 1, xor ^ nums[index])

    bt(0, k, 0)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e7))
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    res = maximizeXor(nums, k)
    print(res)
