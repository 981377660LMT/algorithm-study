"""
题意:给定一个只含有K,E,Y字符的字符串,
计算在最多K次邻位交换后形成的所有可能的字符串集合的大小
len(s)<=30 k<=1e9

!两个同构串的最小邻位交换次数:贪心地从左向右匹配,找最近的换

1. 暴力枚举3^30 不可取
2. 线性dp(index,move,k,e) 表示当前位置,当前交换次数,当前k,e的个数
dfs(index,remain,k,e,y) = dfs(index+1,remain-nextK,k+1,e,y) + ...
时间复杂度O(n^6)
"""

from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    @lru_cache(None)
    def dfs(index: int, remain: int, K: int, E: int, Y: int) -> int:
        if remain < 0:
            return 0
        if index == n:
            return 1 if remain >= 0 else 0

        # !O(n)寻找下一个K/E/Y的位置
        curK, curE, curY, curS = K, E, Y, ""
        for char in s:
            if char == "K":
                if curK == 0:
                    curS += "K"
                else:
                    curK -= 1
            elif char == "E":
                if curE == 0:
                    curS += "E"
                else:
                    curE -= 1
            elif char == "Y":
                if curY == 0:
                    curS += "Y"
                else:
                    curY -= 1

        nextK, nextE, nextY = -1, -1, -1
        for i in range(len(curS)):
            if curS[i] == "K" and nextK == -1:
                nextK = i
            elif curS[i] == "E" and nextE == -1:
                nextE = i
            elif curS[i] == "Y" and nextY == -1:
                nextY = i

        res = 0
        if nextK != -1:
            res += dfs(index + 1, remain - nextK, K + 1, E, Y)
        if nextE != -1:
            res += dfs(index + 1, remain - nextE, K, E + 1, Y)
        if nextY != -1:
            res += dfs(index + 1, remain - nextY, K, E, Y + 1)
        return res

    s = input()
    k = int(input())
    n = len(s)
    res = dfs(0, min(k, n * n), 0, 0, 0)  # !当前位置,剩余可用的交换次数,当前k,e,y的个数
    print(res)
