import sys
from typing import List

input = sys.stdin.readline

N = int(input())
if N & 1:
    print()
    exit()


def dfs(path: List[str], left: int, right: int) -> None:
    if left > right:
        return

    if left == right == 0:
        print(''.join(path))
        return

    if left > 0:
        path.append('(')
        dfs(path, left - 1, right)
        path.pop()

    if right > 0:
        path.append(')')
        dfs(path, left, right - 1)
        path.pop()


dfs([], N // 2, N // 2)

