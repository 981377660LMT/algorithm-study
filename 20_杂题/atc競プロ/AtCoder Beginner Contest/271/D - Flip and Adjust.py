"""
有n个卡牌,正面和反面都有一个权值,
你可以选择一个面朝上,最后使得n张牌朝上的总和为m,输出任意一组方案数。

dp复原
"""
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 上に向けられた面に書かれた整数の総和がちょうど S となるようにカードを置くことができるか判定し、可能ならそのようなカードの置き方の一例を求めてください。
# !H:表 Head T:裏 Tail
if __name__ == "__main__":
    n, s = map(int, input().split())
    goods = [tuple(map(int, input().split())) for _ in range(n)]

    dp = [[False] * (s + 1) for _ in range(n + 1)]
    dp[0][0] = True
    for i, (a, b) in enumerate(goods):
        for cap in range(s + 1):
            if dp[i][cap]:
                if cap + a <= s:
                    dp[i + 1][cap + a] = True
                if cap + b <= s:
                    dp[i + 1][cap + b] = True

    if not dp[n][s]:
        print("No")
        exit(0)

    print("Yes")

    res = []
    cur = s
    for i in range(n - 1, -1, -1):
        if cur - goods[i][0] >= 0 and dp[i][cur - goods[i][0]]:
            res.append("H")
            cur -= goods[i][0]
        else:
            res.append("T")
            cur -= goods[i][1]

    print("".join(res[::-1]))
