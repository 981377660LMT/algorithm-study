import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    res = 0
    bit = 1
    for i in range(60):
        count = 0
        odd = 0
        for j in range(N):
            if A[j] & (1 << i):
                odd = ~odd
            if odd:
                count += 1
        for j in range(N):
            res += bit * count
            if A[j] & (1 << i):
                count = N - j - count
        bit *= 2

    print(res - sum(A))
