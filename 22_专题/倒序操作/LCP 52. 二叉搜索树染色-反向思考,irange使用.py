# 逆向思维
from typing import List, Optional
from sortedcontainers import SortedList


class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


# ops[i] = [type, x, y] 表示第 i 次操作为：
# type 等于 0 时，将节点值范围在 [x, y] 的节点均染蓝
# type 等于 1 时，将节点值范围在 [x, y] 的节点均染红
# 请返回完成所有染色后，该二叉树中红色节点的数量。
# 将所有节点放进SortedList中,倒序遍历所有染色操作,并且记录其每次染的颜色即可
class Solution:
    def getNumber(self, root: Optional[TreeNode], ops: List[List[int]]) -> int:
        def dfs(cur: Optional['TreeNode']):
            if not cur:
                return

            sortedList.add(cur.val)
            dfs(cur.left)
            dfs(cur.right)

        sortedList = SortedList()
        dfs(root)

        # 反向操作 红的加 每次看完就删掉那一段
        res = 0
        for opt, left, right in reversed(ops):
            nums = list(sortedList.irange(left, right))  # irange返回值是chain,要转list

            if opt == 1:
                res += len(nums)

            for num in nums:
                sortedList.discard(num)

        return res

