# 给定一棵树，求包含每个点的路径数
# 每个点出现在多少条路径上
# !路径分两种情况，一种是没有父节点参与的，树形 DP 一下就行了；另一种是父节点参与的，个数就是 子树*(n-子树)


from typing import List


def countPath(n: int, adjList: List[List[int]]) -> List[int]:
    """求包含每个点的路径数.路径至少有两个点."""
    res = [0] * n

    def dfs(cur: int, pre: int) -> int:
        count = 0
        size = 1
        for next_ in adjList[cur]:
            if next_ != pre:
                subSize = dfs(next_, cur)
                count += size * subSize  # 以当前节点为根的子树的路径数
                size += subSize
        count += size * (n - size)  # 父亲节点参与的路径数
        res[cur] = count
        return size

    dfs(0, -1)
    return res
