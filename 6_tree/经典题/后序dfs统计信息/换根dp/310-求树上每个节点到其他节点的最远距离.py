from typing import List


def maxDist(n: int, edges: List[List[int]]) -> List[int]:
    """求树上每个节点(作为根节点)到其他节点的最远距离"""

    def dfs1(cur: int, pre: int) -> int:
        """以0为根到其子树的最远距离"""
        for next in adjList[cur]:
            if next == pre:
                continue
            maxCand = dfs1(next, cur) + 1
            if maxCand > down1[cur]:
                down2[cur], down1[cur] = down1[cur], maxCand
                downMaxNeedRoot[cur] = next
            elif maxCand > down2[cur]:
                down2[cur] = maxCand
        return down1[cur]

    def dfs2(cur: int, pre: int) -> None:
        """以0为根 节点i子树以外的节点到i的最远距离"""
        for next in adjList[cur]:
            if next == pre:
                continue
            if downMaxNeedRoot[cur] == next:
                up[next] = max(up[cur], down2[cur]) + 1
            else:
                up[next] = max(up[cur], down1[cur]) + 1
            dfs2(next, cur)

    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    # 分别记录`向下(子节点)`的最大值和次大值dp1
    down1, down2 = [0] * n, [0] * n
    # `向下`取最大值时必须经过的结点
    downMaxNeedRoot = [0] * n
    # 记录`节点向上(父结点)`的最大距离dp2
    up = [0] * n

    dfs1(0, -1)
    dfs2(0, -1)

    res = [max(up[i], down1[i]) for i in range(n)]
    return res


if __name__ == "__main__":
    print(maxDist(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))
