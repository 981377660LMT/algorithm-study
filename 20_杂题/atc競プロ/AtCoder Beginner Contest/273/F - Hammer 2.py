# F - Hammer 2
# 锤子

# 从原点开始，到目标点target
# 求最少的路程。
# 有n个锤子，n个墙
# h[i] 第i个锤子的位置
# w[i] 第i面墙的位置
# 第i面墙只能被第i个锤子打碎

# !dp是优雅的暴力
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, x = map(int, input().split())
    wall = list(map(int, input().split()))
    hammer = list(map(int, input().split()))

    # 区间dp 需要离散化
    # dp(left,right,flag) 表示(l,r)这个区间都已经遍历过，0表示当前在左端点，1表示当前在右端点
    # !转移时，要考虑x-1和y+1是否是墙，如果是墙，是否拿到了锤子。


# TODO
