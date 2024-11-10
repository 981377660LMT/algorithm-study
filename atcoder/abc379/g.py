import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    N = int(input())
    S = input()

    total = 0
    currentSum = 0

    for i in range(N):
        digit = int(S[i])
        currentSum = currentSum * 10 + digit * (i + 1)
        total += currentSum

    print(total)
