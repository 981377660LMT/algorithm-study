# E - Permute K times 2
# !给定置换P，进行k次操作。
# 每次操作，P[i]=P[P[i]]，问k次操作后的置换P。
# !注意操作是同时进行的.
#
# 找到置换环，每个点移动pow(2,k,len(环大小))次.

from typing import List


def collectCycle(nexts: List[int], start: int) -> List[int]:
    """置换环找环.nexts数组中元素各不相同."""
    cycle = []
    cur = start
    while True:
        cycle.append(cur)
        cur = nexts[cur]
        if cur == start:
            break
    return cycle


if __name__ == "__main__":
    N, K = map(int, input().split())
    P = list(map(int, input().split()))
    for i in range(len(P)):
        P[i] -= 1

    visited = [False] * N
    res = [0] * N
    for i in range(N):
        if visited[i]:
            continue
        cycle = collectCycle(P, i)
        for v in cycle:
            visited[v] = True

        step = pow(2, K, len(cycle))
        for j in range(len(cycle)):
            res[cycle[j]] = cycle[(j + step) % len(cycle)] + 1

    print(" ".join(map(str, res)))
