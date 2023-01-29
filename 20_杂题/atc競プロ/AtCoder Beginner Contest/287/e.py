import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 英小文字からなる文字列が
# N 個与えられます。
# i(i=1,2,…,N) 番目のものを
# S
# i
# ​
#   と表します。

# 二つの文字列
# x,y に対し、以下の条件を全て満たす最大の整数
# n を
# LCP(x,y) と表します。

# x,y の長さはいずれも
# n 以上
# 1 以上
# n 以下の全ての整数
# i に対し、
# x の
# i 文字目と
# y の
# i 文字目が等しい
# 全ての
# i=1,2,…,N に対し、以下の値を求めてください。

# i
# 
# =j
# max
# ​
#  LCP(S
# i
# ​
#  ,S
# j
# ​
#  )


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    wordsWithIndex = [(i, word) for i, word in enumerate(words)]
    wordsWithIndex.sort(key=lambda x: x[1])
    res = [0] * n
    for i in range(n):
        cur = wordsWithIndex[i][1]
        if i > 0:
            pre = wordsWithIndex[i - 1][1]
            count = 0
            for a, b in zip(cur, pre):
                if a == b:
                    count += 1
                else:
                    break
            res[wordsWithIndex[i][0]] = max(res[wordsWithIndex[i][0]], count)
        if i < n - 1:
            nxt = wordsWithIndex[i + 1][1]
            count = 0
            for a, b in zip(cur, nxt):
                if a == b:
                    count += 1
                else:
                    break
            res[wordsWithIndex[i][0]] = max(res[wordsWithIndex[i][0]], count)
    print(*res, sep="\n")
