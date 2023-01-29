# 回文+移动就会联想到循环节啊


def stAndTSPalindrome(s: str, k: int) -> bool:
    """
    给定一个字符串s和一个正整数k,
    判断是否存在一个长为k的字符串t,使得s+t是回文串,且t+s是回文串

    直接 K %= N*4,然后头尾填充 逆序S 和 S,得到T;
    然后直接判断 S+T 和 T+S 是不是回文
    """
    n = len(s)
    k %= n * 4
    res = ""
    cur = s
    while len(res) < k:
        cur = cur[::-1]
        res += cur
    res = res[:k]
    return res + s == (res + s)[::-1] and s + res == (s + res)[::-1]


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    T = int(input())
    for _ in range(T):
        n, k = map(int, input().split())
        s = input()
        print("Yes" if stAndTSPalindrome(s, k) else "No")
