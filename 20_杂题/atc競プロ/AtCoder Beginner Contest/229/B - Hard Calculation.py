# 给两个数A, B，问A+B的计算过程中是否发生了进位


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    A, B = map(int, input().split())
    while A and B:
        div1, mod1, div2, mod2 = A // 10, A % 10, B // 10, B % 10
        if mod1 + mod2 >= 10:
            print("Hard")
            exit(0)
        A, B = div1, div2
    print("Easy")
