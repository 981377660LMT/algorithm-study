from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    intervals = [sorted(map(int, input().split())) for _ in range(N)]
    intervals.sort(key=lambda x: (x[0], -x[1]))
    curStart, curEnd = intervals[0]
    for i in range(1, len(intervals)):
        preStart, preEnd = intervals[i]
        if curEnd > preStart and curEnd < preEnd:
            print("Yes")
            exit(0)
        curStart, curEnd = preStart, preEnd

    print("No")
