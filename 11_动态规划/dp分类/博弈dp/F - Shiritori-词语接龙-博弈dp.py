# shiritori 词语接龙/成语接龙
# 不能用用过的词 且开头必须和上一个词的结尾相同(如果有的话)
# !问先手是否必胜

from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(words: List[str]) -> bool:
    @lru_cache(None)
    def dfs(visited: int, pre: str) -> bool:
        for i in range(n):
            if visited & (1 << i) == 0 and (pre == "" or words[i][0] == pre):
                if not dfs(visited | (1 << i), words[i][-1]):
                    return True
        return False

    return dfs(0, "")


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    print("First" if solve(words) else "Second")
