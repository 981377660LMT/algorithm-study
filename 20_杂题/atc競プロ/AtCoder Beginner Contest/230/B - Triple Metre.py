# 判断S是否为'oxx'*int(1e5) 的子串
# !时间复杂度O(len(S))

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    s = input()
    for cycle in ("oxx", "xxo", "xox"):
        flag = True
        for i in range(len(s)):
            if s[i] != cycle[i % 3]:
                flag = False
                break
        if flag:
            print("Yes")
            exit(0)
    print("No")
