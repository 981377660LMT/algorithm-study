# 树上的LCS
# 树上每条边有一个字符
# 求树上的一条路径，使得路径上的字符序列与给定的字符串target的最长公共子序列最长.
# n<=2000
# s<=2000

# !每个顶点作为起点时的路径
from collections import defaultdict
from typing import List, Tuple
from Rerooting import Rerooting


def lcsOnTree(n: int, edges: List[Tuple[int, int, str]], target: str) -> int:
    E = List[int]

    def e(root: int) -> E:
        return [0] * (len(target) + 1)

    def op(childRes1: E, childRes2: E) -> E:
        res = childRes1[:]
        for i in range(len(res)):
            if res[i] < childRes2[i]:
                res[i] = childRes2[i]
        return res

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        from_, to = (cur, parent) if direction == 0 else (parent, cur)
        c = weights[from_][to]
        res = fromRes[:]
        for i in range(len(s) - 1, -1, -1):
            if s[i] == c:
                if res[i + 1] < res[i] + 1:
                    res[i + 1] = res[i] + 1
        for i in range(len(s)):
            if res[i + 1] < res[i]:
                res[i + 1] = res[i]
        return res

    R = Rerooting(n)
    weights = [defaultdict(str) for _ in range(n)]
    for u, v, c in edges:
        R.addEdge(u, v)
        weights[u][v] = c
        weights[v][u] = c

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    return max(v[-1] for v in dp)


if __name__ == "__main__":
    n = int(input())
    s = input()
    edges = []
    for _ in range(n - 1):
        u, v, c = input().split()
        edges.append((int(u) - 1, int(v) - 1, c))

    print(lcsOnTree(n, edges, s))
