# 进制转换
# 任意进制字符串相互转换的通用解法，注意特判0
# 124. 数的进制转换
# 编写一个程序，可以实现将一个数字由一个进制转换为另一个进制。
# 这里有 62 个不同数位 {0−9,A−Z,a−z}。


import string

allChar = string.digits + string.ascii_lowercase + string.ascii_uppercase
digitToChar = {i: char for i, char in enumerate(allChar)}
charToDigit = {char: i for i, char in enumerate(allChar)}


def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
    assert 2 <= radix <= 62
    if num < 0:
        return "-" + toString(-num, radix)
    if num == 0:
        return "0"
    res = []
    while num:
        div, mod = divmod(num, radix)
        res.append(digitToChar[mod])
        num = div
    return "".join(res)[::-1] or "0"


def parseInt(string: str, radix: int) -> int:
    """将字符串转换为指定进制的数字"""
    return int(string, base=radix)


####################################################################################


def convert(num: str, rawRadix: int, targetRadix: int) -> str:
    """先将原始进制的字符串转换为10进制大数,然后再将这个数转换为目标进制的字符串"""
    assert 2 <= rawRadix, targetRadix <= 62
    if num == "0":
        return "0"
    # 原始进制转10进制
    decimal = 0
    base = 1
    for i in range(len(num) - 1, -1, -1):
        char = num[i]
        decimal += base * charToDigit[char]
        base *= rawRadix
    # 10进制转目标进制
    res = []
    while decimal:
        div, mod = divmod(decimal, targetRadix)
        res.append(digitToChar[mod])
        decimal = div
    return "".join(res)[::-1] or "0"


if __name__ == "__main__":
    n, k = input().split()
    k = int(k)
    for _ in range(k):
        n = convert(n, 8, 9)
        n = n.replace("8", "5")
    print(n)
