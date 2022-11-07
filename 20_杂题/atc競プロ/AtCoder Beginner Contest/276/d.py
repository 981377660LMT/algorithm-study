import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 正整数列 A=(a
# 1
# ​
#  ,a
# 2
# ​
#  ,…,a
# N
# ​
#  ) が与えられます。
# あなたは以下の操作のうち 1 つを選んで行うことを 0 回以上何度でも繰り返せます。

# 1≤i≤N かつ a
# i
# ​
#   が 2 の倍数であるような整数 i を選び、a
# i
# ​
#   を
# 2
# a
# i
# ​

# ​
#   に置き換える
# 1≤i≤N かつ a
# i
# ​
#   が 3 の倍数であるような整数 i を選び、a
# i
# ​
#   を
# 3
# a
# i
# ​

# ​
#   に置き換える
# あなたの目標は A が a
# 1
# ​
#  =a
# 2
# ​
#  =…=a
# N
# ​
#   を満たす状態にすることです。
# 目標を達成するために必要な操作の回数の最小値を求めてください。ただし、どのように操作を行っても目標を達成できない場合、代わりに -1 と出力してください。

# 只能有这些因子
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    fac2 = [0] * n
    fac3 = [0] * n
    for i in range(n):
        while nums[i] % 2 == 0:
            fac2[i] += 1
            nums[i] //= 2
        while nums[i] % 3 == 0:
            fac3[i] += 1
            nums[i] //= 3

    S = set(nums)
    if len(S) != 1:
        print(-1)
        exit(0)

    min2, min3 = min(fac2), min(fac3)
    print(sum(fac2) + sum(fac3) - (min2 + min3) * n)
