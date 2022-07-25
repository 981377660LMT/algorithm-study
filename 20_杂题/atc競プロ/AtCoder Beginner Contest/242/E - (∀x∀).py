# 给定字符串S由大写字母组成 ,问有多少个回文串X字典序小于或等于S ?
# T<=2.5e5


# 找左半对称串 考虑10进制
# 1234
# 答案为 [0000,...,1221]
# 1202
# 答案为 [0000,...,1221] 去掉 1221

# !1.回文由左半串唯一确定
# !2. 26进制计数
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def cal(s: str) -> int:
    """26进制下多少个回文串X字典序小于或等于S"""
    n = len(s)
    left = s[: (n + 1) // 2]
    right = left[::-1]
    if n & 1:
        right = right[1:]
    res = 0
    for char in left:
        res = res * 26 + ord(char) - ord("A")  # 每个位严格小于的个数
        res %= MOD
    if left + right <= s:
        res += 1
        res %= MOD
    return res


def main() -> None:
    T = int(input())
    for _ in range(T):
        n = int(input())
        s = input()
        print(cal(s))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
