import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 個の整数の
# 2 つ組
# (A
# 1
# ​
#  ,B
# 1
# ​
#  ),(A
# 2
# ​
#  ,B
# 2
# ​
#  ),…,(A
# N
# ​
#  ,B
# N
# ​
#  ) が与えられます。 各
# i=1,2,…,N について、
# A
# i
# ​
#  +B
# i
# ​
#   を出力してください。
if __name__ == "__main__":
    n = int(input())
    for _ in range(n):
        a, b = map(int, input().split())
        print(a + b)
