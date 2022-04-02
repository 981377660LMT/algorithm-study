# 给定一个整数 n ，返回 可表示为两个 n 位整数乘积的 最大回文整数 。因为答案可能非常大，所以返回它对 1337 取余 。
from itertools import chain
from typing import Generator, Iterable


MOD = 1337

# 1 <= n <= 8
# 回文串最多17位


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


class Solution:
    def largestPalindrome(self, n: int) -> int:
        iter = genPalindromeByLength(1, 2 * n, isReversed=True)
        lower, upper = 10 ** (n - 1), 10 ** n - 1
        for palindrome in iter:
            # sqrt = int(palindrome ** 0.5)
            # print(sqrt - lower, upper - sqrt)
            # 这里用大的范围，是因为从后往前找，[lower,sqrt]里的数多 [sqrt,upper]里的数少
            for factor in range(upper, int(palindrome ** 0.5), -1):
                # for factor in range(lower, int(palindrome ** 0.5) + 1):
                if palindrome % factor == 0 and upper >= palindrome // factor >= lower:
                    return palindrome % MOD
        return -1


print(Solution().largestPalindrome(1))
print(Solution().largestPalindrome(2))
print(Solution().largestPalindrome(6))
# 两个1位
# 1 1 2
# 2 3 4
# 3 5 6

