from collections import Counter, defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# AtCoder社ではカードを使った
# 1 人ゲームが流行っています。
# ゲームで使う各カードには、英小文字
# 1 文字または @ の文字が書かれており、いずれのカードも十分多く存在します。
# ゲームは以下の手順で行います。

# カードを同じ枚数ずつ
# 2 列に並べる。
# @ のカードを、それぞれ a, t, c, o, d, e, r のいずれかのカードと置き換える。
# 2 つの列が一致していれば勝ち。そうでなければ負け。
# このゲームに勝ちたいあなたは、次のようなイカサマをすることにしました。

# 手順
# 1 以降の好きなタイミングで、列内のカードを自由に並び替えてよい。
# 手順
# 1 で並べられた
# 2 つの列を表す
# 2 つの文字列
# S,T が与えられるので、イカサマをしてもよいときゲームに勝てるか判定してください。
if __name__ == "__main__":
    s = input()
    t = input()
    cur, target = Counter(), Counter()
    at1 = 0
    at2 = 0
    for i in s:
        if i == "@":
            at1 += 1
        else:
            cur[i] += 1

    for i in t:
        if i == "@":
            at2 += 1
        else:
            target[i] += 1

    diff1 = target - cur
    diff2 = cur - target
    diff = diff1 + diff2

    keys = diff.keys()
    if all(k in "@atcoder" for k in keys) and sum(diff.values()) <= at1 + at2:
        print("Yes")
    else:
        print("No")
