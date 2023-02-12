from itertools import combinations
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 1 以上
# N 以下の整数からなる集合が
# M 個あり、順に
# S
# 1
# ​
#  ,S
# 2
# ​
#  ,…,S
# M
# ​
#   と呼びます。
# S
# i
# ​
#   は
# C
# i
# ​
#   個の整数
# a
# i,1
# ​
#  ,a
# i,2
# ​
#  ,…,a
# i,C
# i
# ​

# ​
#   からなります。

# M 個の集合から
# 1 個以上の集合を選ぶ方法は
# 2
# M
#  −1 通りあります。
# このうち、次の条件を満たす選び方は何通りありますか？

# 1≤x≤N を満たす全ての整数
# x に対して、選んだ集合の中に
# x を含む集合が少なくとも
# 1 個存在する。

if __name__ == "__main__":
    n, m = map(int, input().split())
    groups = []
    for _ in range(m):
        size = int(input())
        nums = list(map(int, input().split()))
        groups.append(set(nums))

    res = 0
    for states in range(1, 1 << m):
        cur = set()
        for i in range(m):
            if states >> i & 1:
                cur |= groups[i]
        res += len(cur) == n
    print(res)
