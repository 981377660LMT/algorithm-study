# 一个 k 镜像数字 指的是一个在十进制和 k 进制下从前往后读和从后往前读都一样的 没有前导 0 的 正 整数。
# 2 <= k <= 9
# 1 <= n <= 30
from typing import Generator
from itertools import chain


def genPalindromeByLength(minLen: int, maxLen: int, isReversed=False) -> Generator[int, None, None]:
    """返回minLength<=长度<=maxLength的回文数字"""

    def inner(length: int, isReversed=False) -> Generator[int, None, None]:
        """返回长度为length的回文数字"""
        # 长为3，4的回文都是从10开始的，所以只需要构造10-99的回文即可
        start = 10 ** ((length - 1) >> 1)
        end = start * 10 - 1

        for half in reversed(range(start, end + 1)) if isReversed else range(start, end + 1):
            if length & 1:
                yield (int(str(half)[:-1] + str(half)[::-1]))
            else:
                yield (int(str(half) + str(half)[::-1]))

    return chain.from_iterable(
        inner(len_, isReversed)
        for len_ in (
            reversed(range(minLen, maxLen + 1)) if isReversed else range(minLen, maxLen + 1)
        )
    )


def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
    assert radix >= 2

    if num < 0:
        return '-' + toString(-num, radix)

    if num == 0:
        return '0'

    res = []
    while num:
        div, mod = divmod(num, radix)
        res.append(str(mod))
        num = div
    return ''.join(res)[::-1]


class Solution:
    def kMirror(self, k: int, n: int) -> int:
        res = []
        iter = genPalindromeByLength(1, int(1e20))
        while len(res) < n:
            palindrome = next(iter)
            cand = toString(palindrome, k)
            if cand == cand[::-1]:
                res.append(palindrome)
        return sum(res)


print(Solution().kMirror(k=3, n=7))
# 输出：499
# 解释：
# 7 个最小的 3 镜像数字和它们的三进制表示如下：
#   十进制       三进制
#     1          1
#     2          2
#     4          11
#     8          22
#     121        11111
#     151        12121
#     212        21212
# 它们的和为 1 + 2 + 4 + 8 + 121 + 151 + 212 = 499 。
