from collections import defaultdict
from typing import DefaultDict, Generator, List, Optional, Set
from TreeManager import TreeManager


def nLCA(tree: DefaultDict[int, Set[int]], values: List[int]) -> int:
    """子树结点互不相同 预处理O(nlogn)+归并查询O(logn*logk)"""

    def merge(left: int, right: int) -> int:
        """归并查询 values[left:right+1] 的LCA"""
        if left == right:
            return values[left]
        if left + 1 == right:
            return treeManager.queryLCA(values[left], values[right])
        mid = (left + right) // 2
        leftLCA = merge(left, mid)
        rightLCA = merge(mid + 1, right)
        return treeManager.queryLCA(leftLCA, rightLCA)

    treeManager = TreeManager(len(tree), tree, root=0, useLCA=True)
    return merge(0, len(values) - 1)


class Tree:
    def __init__(self, val, children: Optional[List["Tree"]] = None) -> None:
        self.val = val
        self.children = children if children is not None else []


def nLCA2(root: Tree, values: List[int]) -> Tree:
    """子树结点互不相同 O(n)"""

    def dfs(root: Optional[Tree]) -> Generator[Tree, None, None]:
        """后序dfs和从下往上拓扑排序 都是等价的
        看哪个点`最先`为n
        """
        if not root:
            return

        counter[root.val] += int(root.val in needs)

        for child in root.children:
            yield from dfs(child)
            counter[root.val] += counter[child.val]

        if counter[root.val] == len(values):
            yield root

    needs = set(values)
    counter = defaultdict(int)
    return next(dfs(root))  # 大海捞针适合生成器写法


if __name__ == "__main__":
    adjMap = defaultdict(set)
    for i in range(6):
        adjMap[i].add(i + 1)
        adjMap[i + 1].add(i)
    print(nLCA(adjMap, [3, 4, 5]))
