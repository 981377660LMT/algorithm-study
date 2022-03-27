from typing import Generator
from itertools import chain


def genPalindrome(length: int) -> Generator[int, None, None]:
    # 长为3，4的回文都是从10开始的，所以只需要构造10-99的回文即可
    start = 10 ** ((length - 1) >> 1)
    end = start * 10 - 1

    for half in range(start, end + 1):
        if length & 1:
            yield int(str(half)[:-1] + str(half)[::-1])
        else:
            yield int(str(half) + str(half)[::-1])


def getPalindromeWithMaxLength(maxLength: int) -> chain[int]:
    return chain.from_iterable(genPalindrome(length) for length in range(1, maxLength + 1))
