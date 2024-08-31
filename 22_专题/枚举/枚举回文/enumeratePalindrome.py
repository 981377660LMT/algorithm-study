"""回文生成器

Api:
enumeratePalindrome(min: int, max: int) -> Generator[int, None, None]
emumeratePalindromeByLength(minLength: int, maxLength: int, reverse=False) -> Generator[str, None, None]
getPalindromeByHalf(half: Union[str, int], even=True) -> str
countPalindrome(length: int) -> int
getKthPalindrome(length: int, k: int) -> Optional[str]
nextPalindrome(x: str) -> str
"""

from typing import Generator, Optional, Union


def enumeratePalindrome(min: int, max: int) -> Generator[int, None, None]:
    """
    遍历[min,max]闭区间内的回文数.
    """
    if min > max:
        return
    minLength = len(str(min))
    base = pow(10, (minLength - 1) >> 1)
    while True:
        # 生成奇数长度回文数，例如 base = 10，生成的范围是 101 ~ 999
        for i in range(base, base * 10):
            x = i
            t = i // 10
            while t > 0:
                x = x * 10 + t % 10
                t //= 10
            if x > max:
                return
            if x >= min:
                yield x

        # 生成偶数长度回文数，例如 base = 10，生成的范围是 1001 ~ 9999
        for i in range(base, base * 10):
            x = i
            t = i
            while t > 0:
                x = x * 10 + t % 10
                t //= 10
            if x > max:
                return
            if x >= min:
                yield x

        base *= 10


def emumeratePalindromeByLength(
    minLength: int, maxLength: int, reverse=False
) -> Generator[str, None, None]:
    """
    遍历长度在 `[minLength, maxLength]` 之间的回文数字字符串.
    maxLength <= 12.
    """
    if minLength > maxLength:
        return
    if reverse:
        for length in reversed(range(minLength, maxLength + 1)):
            start = 10 ** ((length - 1) >> 1)
            end = start * 10 - 1
            for half in reversed(range(start, end + 1)):
                s = str(half)
                if length & 1:
                    yield f"{s}{s[:-1][::-1]}"
                else:
                    yield f"{s}{s[::-1]}"
    else:
        for length in range(minLength, maxLength + 1):
            start = 10 ** ((length - 1) >> 1)
            end = start * 10 - 1
            for half in range(start, end + 1):
                s = str(half)
                if length & 1:
                    yield f"{s}{s[:-1][::-1]}"
                else:
                    yield f"{s}{s[::-1]}"


def getPalindromeByHalf(half: Union[str, int], even=True) -> str:
    """给定回文的一半,返回偶数长度/奇数长度的回文字符串."""
    s = str(half)
    if even:
        return f"{s}{s[::-1]}"
    return f"{s}{s[:-1][::-1]}"


def countPalindrome(length: int) -> int:
    """返回长度为length的回文数个数."""
    if length <= 0:
        return 0
    start = pow(10, ((length - 1) >> 1))
    return (start * 10 - 1) - start + 1


def getKthPalindrome(length: int, k: int) -> Optional[str]:
    """返回长度为length的第k个回文数,k>=1."""
    if length <= 0:
        return None
    start = pow(10, ((length - 1) >> 1))
    count = (start * 10 - 1) - start + 1
    if k > count:
        return None
    half = start + k - 1
    s = str(half)
    if length & 1:
        return f"{s}{s[:-1][::-1]}"
    return f"{s}{s[::-1]}"


def nextPalindrome(x: str) -> str:
    """返回比x大的下一个回文数."""
    if x == "9" * len(x):
        return "1" + "0" * (len(x) - 1) + "1"
    if len(x) & 1:
        half = str(int(x[: len(x) // 2 + 1]) + 1)
        return half + half[:-1][::-1]
    else:
        half = str(int(x[: len(x) // 2]) + 1)
        return half + half[::-1]


if __name__ == "__main__":
    count = 0
    for p in emumeratePalindromeByLength(1, 12):
        count += 1
    print(count)
    print(getPalindromeByHalf(123, False))
    print(getKthPalindrome(2, 9))
    print(nextPalindrome("9999"))
    print(list(enumeratePalindrome(1, 100)))
