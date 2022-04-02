from typing import Generator, Iterable
from itertools import chain


def genPalindromeByLength(minLen: int, maxLen: int, isReversed=False) -> Iterable[int]:
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


def isPrime(n: int) -> bool:
    return n >= 2 and all(n % i for i in range(2, int(n ** 0.5) + 1))


class Solution:
    def primePalindrome(self, n: int) -> int:
        """
        求出大于或等于 N 的最小回文素数。
        1 <= N <= 10^8
        """

        for cand in genPalindromeByLength(1, 9):
            if cand < n:
                continue
            if isPrime(cand):
                return cand

        return -1


print(Solution().primePalindrome(13))
# 输出：101
