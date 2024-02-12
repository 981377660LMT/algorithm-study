from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 個の整数
# A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#   が与えられます。
# 以下の条件を満たす整数の組
# (i,j,k) の個数を求めてください。

# 1≤i,j,k≤N
# A
# i
# ​
#  ×A
# j
# ​
#  =A
# k
# ​


if __name__ == "__main__":
    n = int(input())
    nums = [int(input()) for _ in range(n)]
    mp = dict()
    for num in nums:
        mp[num] = mp.get(num, 0) + 1

    res = 0
    for i in range(n):
        for j in range(n):
            res += mp.get(nums[i] * nums[j], 0)
    print(res)
