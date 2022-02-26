# n<=1000
# 0..n 的所有数按升序组成的序列，
# 我们要进行一些筛选，每次我们丢弃去当前所有数字中第奇数位个的数。
# 重复这一过程直到最后剩下一个数。
# 请求出最后剩下的数字。


# 注意到每次移除二进制下最右边为0的位置
# 最后剩一个数肯定是0到n中二进制下1最多的那个数


def main(n: int) -> None:
    res = 1
    while res <= n:
        res <<= 1
    print((res >> 1) - 1)


while True:
    try:
        n = int(input())
        main(n)
    except EOFError:
        break

