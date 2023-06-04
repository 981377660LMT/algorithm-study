"""
给定两个长为n的字符串s和t
问s和t的所有轮转的子串中 s的轮转子串有多少个字典序 <= t的轮转子串

技巧:
需要一起比较s和t的所有轮转字串的字典序
!构造一个新的字符串 s+s+'#'+t+t+'|'
(注意题目要的是小于等于, 这样保证两个字符串在比较完长度为n后S后面的#小于T中任意一个字符。)
!后缀数组求出每个串的rank
!然后在t的rank中 用s的每个子串rank二分出t中的pos
"""

from bisect import bisect_left
import sys
from SuffixArray import SuffixArray

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# !f(S,i)表示轮转 字符串S `i次``
# 文字列 X と整数 i に対し、f(X,i) を X に対して以下の操作を i 回行い得られる文字列とします。
# X の先頭の文字を削除し、同じ文字を X の末尾に挿入する。
# !0≤i,j≤N−1 を満たす正整数の組 (i,j) のうち、
# !辞書順で f(S,i) が f(T,j) より小さいか同じであるものの個数を求めてください。

if __name__ == "__main__":
    n = int(input())
    s = input()
    t = input()

    SMALL, BIG = chr(0), chr(0x10FFFF)
    sstt = s + s + SMALL + t + t + BIG  # 保证两个字符串在比较完长度为n后 S+'#' 小于T+'|'中任意一个字符。
    ords = list(map(ord, sstt))
    rank = SuffixArray(ords).rank
    sRank, tRank = rank[:n], rank[2 * n + 1 : 2 * n + 1 + n]
    tRank.sort()

    res = 0
    for cur in sRank:
        res += n - bisect_left(tRank, cur)
    print(res)
