# 无向图中的不传递三元组
# find_nontransitive_triple
# https://maspypy.github.io/library/graph/find_nontransitive_triple.hpp
# https://codeforces.com/contest/967/problem/F
# ab, bc 辺はあるが ac 辺はないような 3 つ組 (a,b,c) を探す。
# なければ None を返す。

from typing import List, Optional, Tuple


def findNontransitiveTriple(n: int, adjList: List[List[int]]) -> Optional[Tuple[int, int, int]]:
    done = [0] * n
    que = []
    for root in range(n):
        if done[root]:
            continue
        que = [root]
        p = 0
        while p < len(que):
            v = que[p]
            p += 1
            done[v] = 2
            s = 0
            for to in adjList[v]:
                if done[to] == 0:
                    done[to] = 1
                    que.append(to)
                elif done[to] == 2:
                    s += 1
            if s == p - 1:
                continue
            # assert p >= 3
            c = v
            a = -1
            b = -1
            for to in adjList[v]:
                done[to] = 0
            for i in range(p - 1):
                x = que[i]
                if done[x] == 2:
                    a = x
                if done[x] == 0:
                    b = x
            # assert a != -1
            # assert b != -1
            return a, b, c

    return (-1, -1, -1)
