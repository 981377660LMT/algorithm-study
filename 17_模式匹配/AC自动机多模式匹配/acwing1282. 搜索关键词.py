# https://www.acwing.com/problem/content/1284/
# https://www.acwing.com/solution/content/50169/

# 给定 n 个长度不超过 50 的由小写英文字母组成的单词，以及一篇长为 m 的文章。
# !请问，其中有多少个单词在文章中出现了。
# !注意：每个单词不论在文章中出现多少次，仅累计 1 次。

from typing import List
from AutoMaton import AhoCorasick

if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    def solve(big: str, smalls: List[str]) -> int:
        ac = AhoCorasick(smalls)
        visited = [False] * len(smalls)
        info = ac.search(big)
        for *_, wid in info:
            visited[wid] = True
        return sum(visited)

    t = int(input())
    for _ in range(t):
        n = int(input())
        smalls = [input() for _ in range(n)]
        big = input()
        print(solve(big, smalls))
