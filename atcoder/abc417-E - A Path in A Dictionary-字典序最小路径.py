# abc417-E - A Path in A Dictionary-字典序最小路径
# https://atcoder.jp/contests/abc417/tasks/abc417_e
# 给定连通无向图 G，有 N 个顶点、M 条边，编号 1…N。
# 给定起点 X、终点 Y。要求一条简单路径 P 从 X 到 Y，使得作为整数序列 P 在字典序上最小。
# 保证答案存在。
# T(≤500) 个测试用例，总 N≤1000，总 M≤5×10⁴。
#
# !https://atcoder.jp/contests/abc417/editorial/13589
# !在 dfs 求路径的基础上，加一个排序即可


from typing import List


def minLexPath(adjList: List[List[int]], start: int, end: int) -> List[int]:
    """字典序最小路径."""
    path = []
    visited = set()

    def dfs(cur: int) -> bool:
        path.append(cur)
        visited.add(cur)
        if cur == end:
            return True
        for next_ in sorted(adjList[cur]):  # !按照字典序遍历
            if next_ not in visited:
                if dfs(next_):
                    return True
        path.pop()
        return False

    dfs(start)

    return path


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    def solve():
        N, M, X, Y = map(int, input().split())
        adjList = [[] for _ in range(N)]
        for _ in range(M):
            u, v = map(int, input().split())
            adjList[u - 1].append(v - 1)
            adjList[v - 1].append(u - 1)

        res = minLexPath(adjList, X - 1, Y - 1)
        print(" ".join(str(x + 1) for x in res))

    T = int(input())
    for _ in range(T):
        solve()
