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

    pre = [[""] * (s + 1) for _ in range(n + 1)]
    dp = [[False] * (s + 1) for _ in range(n + 1)]
    dp[0][0] = True
    for i, (a, b) in enumerate(goods):
        for cap in range(s + 1):
            if dp[i][cap]:
                if cap + a <= s:
                    dp[i + 1][cap + a] = True
                    pre[i + 1][cap + a] = "H"
                if cap + b <= s:
                    dp[i + 1][cap + b] = True
                    pre[i + 1][cap + b] = "T"

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
