# https://atcoder.jp/contests/abc270/tasks/abc270_g
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(P: int, A: int, B: int, S: int, G: int) -> int:
    ...


# TODO

if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        P, A, B, S, G = map(int, input().split())
        print(solve(P, A, B, S, G))
