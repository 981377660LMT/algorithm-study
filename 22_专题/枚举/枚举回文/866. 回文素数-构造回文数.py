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


def isPrime(n: int) -> bool:
    return n >= 2 and all(n % i for i in range(2, int(n ** 0.5) + 1))


class Solution:
    def primePalindrome(self, n: int) -> int:
        """
        求出大于或等于 N 的最小回文素数。
        1 <= N <= 10^8
        """

        for cand in genPalindromeByRange(1, 9):
            if cand < n:
                continue
            if isPrime(cand):
                return cand

        return -1


print(Solution().primePalindrome(13))
# 输出：101
