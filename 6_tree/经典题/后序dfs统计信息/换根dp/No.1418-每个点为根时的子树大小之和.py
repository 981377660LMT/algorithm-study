# https://yukicoder.me/problems/no/1418
# No.1418-每个点为根时的子树大小之和

import sys


sys.setrecursionlimit(int(1e9))

if __name__ == "__main__":
    from Rerooting import Rerooting

    def e(root: int) -> int:
        return 0

    def op(childRes1: int, childRes2: int) -> int:
        return childRes1 + childRes2

    def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        if direction == 0:  # cur -> parent
            return fromRes + subSize[cur]
        return fromRes + (n - subSize[cur])  # parent -> cur

    def dfsForSubSize(cur: int, parent: int) -> int:
        res = 1
        for next in R.adjList[cur]:
            if next != parent:
                res += dfsForSubSize(next, cur)
        subSize[cur] = res
        return res

    n = int(input())
    R = Rerooting(n)
    for _ in range(n - 1):
        u, v = map(int, input().split())
        R.addEdge(u - 1, v - 1)

    subSize = [0] * n
    dfsForSubSize(0, -1)
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    print(sum(dp) + n * n)  # !加上每个根自己的答案
