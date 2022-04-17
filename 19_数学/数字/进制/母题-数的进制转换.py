# 任意进制字符串相互转换的通用解法，注意特判0


# 124. 数的进制转换
# 编写一个程序，可以实现将一个数字由一个进制转换为另一个进制。
# 这里有 62 个不同数位 {0−9,A−Z,a−z}。
import string


allChar = string.digits + string.ascii_uppercase + string.ascii_lowercase
charByDigit = {i: char for i, char in enumerate(allChar)}
digitByChar = {char: i for i, char in enumerate(allChar)}


def convert(num: str, rawRadix: int, targetRadix: int) -> str:
    """先将原始进制的字符串转换为10进制大数,然后再将这个数转换为目标进制的字符串"""
    assert 2 <= rawRadix, targetRadix <= 62

    if num == '0':
        return '0'

    # 原始进制转10进制
    decimal = 0
    base = 1
    for i in range(len(num) - 1, -1, -1):
        char = num[i]
        decimal += base * digitByChar[char]
        base *= rawRadix

    # 10进制转目标进制
    res = []
    while decimal:
        div, mod = divmod(decimal, targetRadix)
        res.append(charByDigit[mod])
        decimal = div
    return ''.join(res)[::-1] or '0'


print(convert('7', 10, 7))

# n = int(input())
# res = []

# for _ in range(n):
#     radix1, radix2, num = input().split()  # 输入进制 输出进制 输入数字
#     radix1, radix2 = int(radix1), int(radix2)
#     res = convert(num, radix1, radix2)
#     print(radix1, num)  # 默认sep为空格,end为'\n'
#     print(radix2, res)
#     print()


# 输入样例：
# 8
# 62 2 abcdefghiz
# 10 16 1234567890123456789012345678901234567890
# 16 35 3A0C92075C0DBF3B8ACBC5F96CE3F0AD2
# 35 23 333YMHOUE8JPLT7OX6K9FYCQ8A


# 输出样例：
# 62 abcdefghiz
# 2 11011100000100010111110010010110011111001001100011010010001

# 10 1234567890123456789012345678901234567890
# 16 3A0C92075C0DBF3B8ACBC5F96CE3F0AD2
