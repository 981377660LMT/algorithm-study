import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    a, b = map(int, input().split())

    # !すぬけ君は、高橋君と青木君のうち少なくとも一方が解けた問題は解け、
    # !2 人とも解けなかった問題は解けませんでした
    print(a | b)
