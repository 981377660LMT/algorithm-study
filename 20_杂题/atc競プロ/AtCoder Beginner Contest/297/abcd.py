import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 正整数
# A,B が与えられます。

# あなたは、
# A=B になるまで以下の操作を繰り返します。

# A,B の大小関係に応じて、次の
# 2 個のうちどちらかの処理を行う。
# A>B ならば、
# A を
# A−B で置き換える。
# A<B ならば、
# B を
# B−A で置き換える。
# A=B になるまで、操作を何回行うか求めてください。ただし、有限回の操作で
# A=B になることが保証されます。
if __name__ == "__main__":
    A, B = map(int, input().split())
    if A == B:
        print(0)
        exit()

    res = 0
    # 取模加速
    while A != B and A % B and B % A:
        if A > B:
            A, B = B, A
        res += B // A
        B %= A
    if A == B:
        print(res)
        exit()
    if A < B:
        A, B = B, A
    res += (A // B) - 1
    print(res)

    # !TODO GCD???
