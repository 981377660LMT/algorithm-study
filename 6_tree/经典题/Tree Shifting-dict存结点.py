# 将子节点向右移动
from collections import defaultdict


class Tree:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def solve(self, root):
        nodemap = defaultdict(list)
        if root:
            nodemap[0].append(root)

        depth = 0
        while nodemap[depth]:
            for node in nodemap[depth]:
                for child in (node.left, node.right):
                    if child:
                        nodemap[depth + 1].append(child)
            depth += 1

        for depth in range(depth - 1, -1, -1):
            row = nodemap[depth + 1]
            for node in reversed(nodemap[depth]):
                node.right = row.pop() if row else None
                node.left = row.pop() if row else None

        return root
