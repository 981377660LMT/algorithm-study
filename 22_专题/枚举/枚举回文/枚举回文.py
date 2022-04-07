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

    yield from chain.from_iterable(
        inner(len_, isReversed)
        for len_ in (
            reversed(range(minLen, maxLen + 1)) if isReversed else range(minLen, maxLen + 1)
        )
    )


############################################################################################
def getPalindromeByHalf(half: str, length: int) -> int:
    """指定回文的一半，返回长为length的回文"""
    if length & 1:
        return int(str(half)[:-1] + str(half)[::-1])
    else:
        return int(str(half) + str(half)[::-1])

