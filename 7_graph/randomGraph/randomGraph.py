# https://maspypy.github.io/library/random/random_graph.hpp


from random import randint, shuffle
from typing import List, Tuple


def randomGraph(n: int, directed: bool, simple: bool) -> List[Tuple[int, int]]:
    """生成随机图的边."""
    edges = []
    cand = []
    for a in range(n):
        for b in range(n):
            if simple and a == b:
                continue
            if not directed and a > b:
                continue
            cand.append((a, b))
    m = randint(0, len(cand))
    s = set()
    for _ in range(m):
        while True:
            i = randint(0, len(cand) - 1)
            if simple and i in s:
                continue
            s.add(i)
            a, b = cand[i]
            edges.append((a, b))
            break
    randomRelabel(n, edges)
    return edges


def randomTree(n: int) -> List[Tuple[int, int]]:
    """生成随机树的边."""
    edges = []
    for i in range(1, n):
        edges.append((randint(0, i - 1), i))
    randomRelabel(n, edges)
    return edges


def randomRelabel(n: int, edges: List[Tuple[int, int]]) -> None:
    shuffle(edges)
    a = list(range(n))
    shuffle(a)
    for i, (u, v) in enumerate(edges):
        edges[i] = (a[u], a[v])


if __name__ == "__main__":
    print(randomGraph(3, False, False))
    print(randomTree(3))
