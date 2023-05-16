# https://atcoder.jp/contests/abc301/tasks/abc301_d


# 给定一个包含01和 ?的字符串s，将?替换成0或 1，
# 使得其表示的二进制是最大的，且不大于n。
# 问这个的最大值。
# 1<=len(s)<=60
# 1<=n<=1e18


# 由于二进制下，任意个数的低位的1
# 加起来都不如一个高位的 1。
# !因此我们就从高位考虑每个 ?，如果替换成1之后不超过 t，就换成 1，否则就换成 0


def bitMask(s: str, n: int) -> int:
    res = 0
    for i, v in enumerate(s):
        if v == "1":
            res |= 1 << (len(s) - i - 1)
    if res > n:
        return -1
    for i, v in enumerate(s):
        if v == "?":
            add = 1 << (len(s) - i - 1)
            if res | add <= n:
                res |= add
    return res


if __name__ == "__main__":
    s = input()
    n = int(input())
    res = bitMask(s, n)
    print(res)
