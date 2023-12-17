"""回文生成器"""

from typing import Generator, Optional, Union


def emumeratePalindrome(
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
                if length & 1:
                    yield f"{half}{str(half)[:-1][::-1]}"
                else:
                    yield f"{half}{str(half)[::-1]}"
    else:
        for length in range(minLength, maxLength + 1):
            start = 10 ** ((length - 1) >> 1)
            end = start * 10 - 1
            for half in range(start, end + 1):
                if length & 1:
                    yield f"{half}{str(half)[:-1][::-1]}"
                else:
                    yield f"{half}{str(half)[::-1]}"


def getPalindromeByHalf(half: Union[str, int], even=True) -> str:
    """给定回文的一半,返回偶数长度/奇数长度的回文字符串."""
    if even:
        return f"{half}{str(half)[::-1]}"
    return f"{half}{str(half)[:-1][::-1]}"


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
    if length & 1:
        return f"{half}{str(half)[:-1][::-1]}"
    return f"{half}{str(half)[::-1]}"


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
    for p in emumeratePalindrome(1, 12):
        count += 1
    print(count)
    print(getPalindromeByHalf(123, False))
    print(getKthPalindrome(2, 9))
    print(nextPalindrome("9999"))
