# !只包含特定字符的最长子串长度

# !如果不合法 就 left = right+1

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

OK = set("ACGT")
if __name__ == "__main__":
    s = input()
    left, res = 0, 0
    for right in range(len(s)):
        if s[right] not in OK:
            left = right + 1
        res = max(res, right - left + 1)
    print(res)
