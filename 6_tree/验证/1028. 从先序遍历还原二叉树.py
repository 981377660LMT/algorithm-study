from typing import Optional

# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# 在遍历中的每个节点处，我们输出 D 条短划线（其中 D 是该节点的深度），然后输出该节点的值。


# 思路:用字典存储节点及其深度，依据`深度`来构建树
# 先序遍历的特性可以知道，上一个遍历到的比自己深度浅1的节点必为自己的父节点
# 根节点的深度为 0
class Solution:
    def recoverFromPreorder(self, traversal: str) -> Optional['TreeNode']:
        # 字典树
        record = {-1: TreeNode(0)}

        # 直接加入或覆盖键key为深度dep的树字典
        def addTree(val: str, depth: int) -> None:
            record[depth] = TreeNode(int(val))
            if not record[depth - 1].left:
                record[depth - 1].left = record[depth]
            else:
                record[depth - 1].right = record[depth]

        val, depth = '', 0
        # 添加哨兵
        for char in traversal + '-':
            if char != '-':
                val += char
            elif val:
                addTree(val, depth)
                val, depth = '', 1
            else:
                depth += 1

        return record[0]


print(Solution().recoverFromPreorder("1-2--3---4-5--6---7"))
print(Solution().recoverFromPreorder("1-2--3--4-5--6--7"))
# 输出：[1,2,5,3,null,6,null,4,null,7]
