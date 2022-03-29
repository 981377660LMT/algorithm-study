# 如果一个正整数自身是回文数，而且它也是一个回文数的平方，那么我们称这个数为超级回文数。
# 返回包含在范围 [L, R] 中的超级回文数的数目。
# L 和 R 是表示 [1, 10^18) 范围的整数的字符串。

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


# 打表要写在外面
palindromes = [*genPalindromeByRange(1, 9)]


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
