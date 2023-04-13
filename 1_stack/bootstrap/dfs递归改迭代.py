# dfs递归改成栈实现
# https://tjkendev.github.io/procon-library/python/graph/dfs.html
from typing import List


def dfs_rec(N, G, s):
    used = [0] * N

    def dfs(v, p):
        used[v] = 1
        ...  # A
        for w in G[v]:
            ...  # B1
            if used[w]:
                continue
            ...  # B2
            r = dfs(w, v)
            ...  # C
        ...  # D
        return ...  # E

    r0 = dfs(s, -1)


# A, B1, B2, C, D, E は上記の関数再帰実装のものと対応する。
def dfs_stack(N, G, s):
    stk = [s]
    used = [0] * N
    R = [None] * N
    it = [0] * N
    while stk:
        v = stk[-1]
        p = stk[-2] if len(stk) > 1 else -1
        if it[v] == 0:
            used[v] = 1
            ...  # A
        else:
            w = G[v][it[v] - 1]
            r = R[w]
            ...  # C

        while it[v] < len(G[v]):
            w = G[v][it[v]]
            it[v] += 1
            ...  # B1
            if used[w]:
                continue
            ...  # B2
            stk.append(w)
            break
        else:
            ...  # D
            R[v] = ...  # E
            stk.pop()
    r0 = R[s]
