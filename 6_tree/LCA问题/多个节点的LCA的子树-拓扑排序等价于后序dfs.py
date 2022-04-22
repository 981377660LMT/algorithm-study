from collections import defaultdict
from typing import Generator, List, Optional


class Tree:
    def __init__(self, val, children: Optional[List['Tree']] = None) -> None:
        self.val = val
        self.children = children if children is not None else []


class Solution:
    def nLCA(self, root: Tree, values: List[int]):
        """子树结点互不相同"""

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

