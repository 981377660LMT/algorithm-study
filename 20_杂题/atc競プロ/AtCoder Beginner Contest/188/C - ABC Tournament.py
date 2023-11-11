import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# tournament 循环赛 求亚军运动员编号
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    dp = [(num, i) for i, num in enumerate(nums)]
    lastLose = 0
    while len(dp) > 1:
        ndp = []
        for i in range(0, len(dp), 2):
            if dp[i][0] > dp[i + 1][0]:
                ndp.append(dp[i])
                lastLose = dp[i + 1][1]
            else:
                ndp.append(dp[i + 1])
                lastLose = dp[i][1]
        dp = ndp

    print(lastLose + 1)
