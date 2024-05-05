import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# (1,2,…,N) を並び替えて得られる数列
# P=(P
# 1
# ​
#  ,P
# 2
# ​
#  ,…,P
# N
# ​
#  ) が与えられます。

# 長さ
# K の正整数列
# (i
# 1
# ​
#  ,i
# 2
# ​
#  ,…,i
# K
# ​
#  ) であって、以下の条件を共に満たすものを良い添字列と呼びます。

# 1≤i
# 1
# ​
#  <i
# 2
# ​
#  <⋯<i
# K
# ​
#  ≤N
# (P
# i
# 1
# ​

# ​
#  ,P
# i
# 2
# ​

# ​
#  ,…,P
# i
# K
# ​

# ​
#  ) はある連続する
# K 個の整数を並び替えることで得られる。
# 厳密には、ある整数
# a が存在して、
# {P
# i
# 1
# ​

# ​
#  ,P
# i
# 2
# ​

# ​
#  ,…,P
# i
# K
# ​

# ​
#  }={a,a+1,…,a+K−1}。
# 全ての良い添字列における
# i
# K
# ​
#  −i
# 1
# ​
#   の最小値を求めてください。 なお、本問題の制約下では良い添字列が必ず
# 1 つ以上存在することが示せます。
# TODO1:
# go snippet window
# !fast set 封装最长连续

# !找到长度为k的子序列，对应a，a+1，a+2，a+k-1，子序列索引差值最小
if __name__ == "__main__":
    N, K = map(int, input().split())
    P = list(map(int, input().split()))
n = len(nums)
res, curSum = INF, 0
for right in range(n):
    curSum += nums[right]
    if right >= k:
        curSum -= nums[right - k]
    if right >= k - 1:
        res = min(res, curSum)
return res
    def check(mid: int) -> bool:
        """长度为mid的窗口中是否存在连续的k个数"""
        res = 0
        for num in nums:
            res += num // mid
        return res >= k

    left, right = K - 1, N - 1
    ok = False
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
            ok = True
        else:
            left = mid + 1
    print(left)
