"""DAG的最小路径覆盖

链是点的集合,这个集合中任意两个元素v、u,要么v能走到u,要么u能走到v。
反链是点的集合,这个集合中任意两点`谁也不能走到谁`。

原题求最大反链
最大反链 = 最小链覆盖(Dilworth 引理)
最小链覆盖就是用最少的链，经过所有的点至少一次
"""

# 给定长度 <=200 的字符串 ，
# 问最多能从里面选出多少个`回文子串`，
# 使得一个子串不是另一个子串的子串

# !子串间的包含关系是一种偏序关系，可以形成一个 DAG
# !求DAG的最小路径覆盖
# !答案为回文子串数量减去二分图最大匹配

from hungarian import Hungarian

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(s: str) -> int:
    n = len(s)
    palindromeSet = set()
    for left in range(n):
        for right in range(left + 1, n + 1):
            cand = s[left:right]
            if cand == cand[::-1]:
                palindromeSet.add(cand)

    palindrome = sorted(palindromeSet, key=len)
    n = len(palindrome)
    H = Hungarian(n, n)
    for r in range(n):
        for c in range(n):
            if r != c and palindrome[r] in palindrome[c]:
                H.addEdge(r, c)

    match = H.work()
    return n - match  # DAG最小路径覆盖


s = input()
print(solve(s))
