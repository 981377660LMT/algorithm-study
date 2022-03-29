# 一个 k 镜像数字 指的是一个在十进制和 k 进制下从前往后读和从后往前读都一样的 没有前导 0 的 正 整数。
# 2 <= k <= 9
# 1 <= n <= 30

from typing import Generator, Iterable
from itertools import chain


def genPalindromeByRange(minLen: int, maxLen: int) -> Iterable[int]:
    """返回 minLength<=长度<=maxLength 的回文数字"""
    return chain.from_iterable(genPalindromeByLength(len_) for len_ in range(minLen, maxLen + 1))


def genPalindromeByLength(length: int) -> Generator[int, None, None]:
    """返回长度为length的回文数字"""
    start = 10 ** ((length - 1) >> 1)
    end = start * 10 - 1
    for half in range(start, end + 1):
        if length & 1:
            yield (int(str(half)[:-1] + str(half)[::-1]))
        else:
            yield (int(str(half) + str(half)[::-1]))


def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
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


# 预处理
palindromes = [*genPalindromeByRange(1, 12)]


class Solution:
    def kMirror(self, k: int, n: int) -> int:
        res = []
        index = 0
        while len(res) < n:
            cand = toString(palindromes[index], k)
            if cand == cand[::-1]:
                res.append(palindromes[index])
            index += 1
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
