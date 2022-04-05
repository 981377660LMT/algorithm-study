# 如果一个正整数自身是回文数，而且它也是一个回文数的平方，那么我们称这个数为超级回文数。
# 返回包含在范围 [L, R] 中的超级回文数的数目。
# L 和 R 是表示 [1, 10^18) 范围的整数的字符串。

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


# 打表要写在外面
palindromes = [*genPalindromeByLength(1, 9)]


class Solution:
    def superpalindromesInRange(self, left: str, right: str) -> int:
        res = []
        lower, upper = int(left), int(right)
        for num in palindromes:
            square = num ** 2
            if square >= lower and square <= upper and square == int(str(square)[::-1]):
                res.append(square)
        return len(res)


print(Solution().superpalindromesInRange(left="4", right="1000"))
# 输出：4
# 解释：
# 4，9，121，以及 484 是超级回文数。
# 注意 676 不是一个超级回文数： 26 * 26 = 676，但是 26 不是回文数。
