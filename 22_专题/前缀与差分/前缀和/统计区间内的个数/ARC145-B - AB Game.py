"""
有一个数字N。Alice和Bob接下来将会玩N场游戏,
第i场游戏共有i个石子,由Alice和Bob轮流操作,Alice 执先手。
每个回合Alice需要从石子里面拿走A的任意正倍数个石子,
Bob需要从石子里面拿走B的任意正倍数个石子。
如果轮到哪个玩家,玩家无法操作,就判这个玩家负。
问当双方都使用最佳策略时,Alice可以赢几场。

n,a,b<=1e18

贪心 玩家要尽可能多的拿数字(一步将军)
a>n时 赢0盘
a<=n时
那此时题目就转化为，在[a,N]中,有多少个大于等于A且mod A小于B的数字。
这种可以前缀和相减做
"""

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, a, b = map(int, input().split())

if a > n:
    print(0)
    exit(0)


def cal(upper: int) -> int:
    """[0,upper]中多少个数模a小于b"""
    div, mod = divmod(upper, a)
    return div * min(a, b) + min(mod, b - 1)


print(cal(n) - cal(a - 1))
