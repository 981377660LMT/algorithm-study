import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    S = input()
    leftCount = [0] * 26
    rightCount = [0] * 26
    for i in range(len(S)):
        rightCount[ord(S[i]) - 65] += 1
    res = 0
    for i in range(len(S)):
        rightCount[ord(S[i]) - 65] -= 1
        for j in range(26):
            res += leftCount[j] * rightCount[j]
        leftCount[ord(S[i]) - 65] += 1

    print(res)
