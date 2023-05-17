# 给定一棵有根树，判断它是否左右对称。
# 对称的树
# https://www.luogu.com.cn/problem/CF1800G
# https://zhuanlan.zhihu.com/p/619412980
# 对所有子树进行哈希，然后从根往下遍历，如果根的所有子树的哈希值都出现偶数次，
# 那么就是对称的，如果有奇数，
# 那么只能有一个，且这个子树要是对称的，那么递归下去就可以了。


from collections import defaultdict
from typing import List, Tuple
from 有根树的同构 import rootedTreeIsomorphism


def symmetree(n: int, edges: List[Tuple[int, int]]) -> int:
    def dfs(cur: int, pre: int) -> bool:
        subHash = defaultdict(int)
        counter = defaultdict(int)
        for next in adjList[cur]:
            if next == pre:
                continue
            hash_ = dp[next]
            subHash[next] = hash_
            counter[hash_] += 1
        if all(v & 1 == 0 for v in counter.values()):
            return True
        if sum(v & 1 == 1 for v in counter.values()) > 1:
            return False
        for next, hash_ in subHash.items():
            if counter[hash_] & 1 == 1:
                return dfs(next, cur)
        raise RuntimeError("unreachable")

    dp = rootedTreeIsomorphism(n, edges)
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    return dfs(0, -1)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    T = int(input())
    for _ in range(T):
        n = int(input())
        edges = []
        for _ in range(n - 1):
            u, v = map(int, input().split())
            u, v = u - 1, v - 1
            edges.append((u, v))
        res = symmetree(n, edges)
        print("YES" if res else "NO")
