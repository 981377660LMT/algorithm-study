import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 長さ
# N の非負整数列
# A が与えられます。
# A から
# K 要素を選んで順序を保ったまま抜き出した列を
# B としたとき、
# MEX(B) としてありえる最大値を求めてください。

# 但し、数列
# X に対して
# MEX(X) は以下の条件を満たす唯一の非負整数
# m として定義されます。

# 0≤i<m を満たす全ての整数
# i が
# X に含まれる。
# m が
# X に含まれない。
if __name__ == "__main__":
    n, k = map(int, input().split())
    a = list(map(int, input().split()))
    a = sorted(set(a))
    a = a[:k]
    S = set(a)
    mex = 0
    while mex in S:
        mex += 1
    print(mex)
