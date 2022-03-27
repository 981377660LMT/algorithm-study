1. 折半构造 intLength<=9
2. 折半找规律 length<=15

```Python
def genPalindromeByLength(length: int) -> Generator[int, None, None]:
    """顺序返回长度为length的所有回文数字"""

    # 长为3，4的回文都是从10开始的，所以只需要构造10-99的回文即可
    start = 10 ** ((length - 1) >> 1)
    end = start * 10 - 1

    for half in range(start, end + 1):
        if length & 1:
            yield (int(str(half)[:-1] + str(half)[::-1]))
        else:
            yield (int(str(half) + str(half)[::-1]))

```
