# 题意:现有n个字符集 {a, b, c, d, . . ..,z}。可以的操作是:
# ·选定第i个字符集。
# ·选用其中的字符，打出L长度的单词。
# 问，至多能打出多少不同的L长度的单词，对998244353 取模。

# !数据范围很明显告诉我们这题是容斥原理

# n 为 3 时
# (只使用第一个键盘)＋(只使用第二个键盘)＋(只使用第三个键盘)
# –(只使用第一个键盘和第二个键盘共有的字母打印)–(只使用第一个键盘和第三个键盘共有的字母打印)-(只使用第二个键盘和第三个键盘共有的字母打印)
# +(使用所用键盘共有的字母打印)。

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n, l = map(int, input().split())
    words = [0] * n
    for i in range(n):
        s = input()
        for char in s:
            words[i] |= 1 << (ord(char) - 97)

    # !枚举用哪些键盘的字母打印
    res = 0
    for state in range(1, 1 << n):
        s = (1 << 26) - 1
        count = 0
        for i in range(n):
            if (state >> i) & 1:
                s &= words[i]
                count += 1
        res += pow(bin(s).count("1"), l, MOD) * (1 if count & 1 else -1)
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
