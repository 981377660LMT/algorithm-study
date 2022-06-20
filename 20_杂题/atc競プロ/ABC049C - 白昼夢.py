# Tが空文字列である状態から始め、
# 以下の操作を好きな回数繰り返すことで S=T とすることができるか判定してください。

# T の末尾に dream dreamer erase eraser のいずれかを追加する。

# len(s)<=1e5

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))


@lru_cache(None)
def dfs(index: int) -> int:
    if index >= n:
        return index == n

    for next in ['dreamer', 'dream', 'eraser', 'erase']:
        if S[index : index + len(next)] == next and dfs(index + len(next)):
            return True
    return False


S = input()
n = len(S)
res = dfs(0)
dfs.cache_clear()
print('YES' if res else 'NO')

