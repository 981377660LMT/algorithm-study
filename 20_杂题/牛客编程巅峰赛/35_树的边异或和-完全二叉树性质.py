from collections import defaultdict
from typing import List

# 完全二叉树
# 树的边异或和
# 建立树的异或路径sum,先求左右节点left/right的公共父节点root，
# （求公共父节点也有原题，用两个变量存是否存在子节点dfs即可）。
# 答案就是(sum[left] ^ sum[root]) ^ (sum[right] ^ sum[root]) ^ root->val

# (sum[left] ^ sum[root]) ^ (sum[right] ^ sum[root]) ^ root->val
# 异或路径


class Solution:
    def tree5(self, preOrder: List[int]) -> int:
        def dfs(root: int) -> None:
            """前序遍历;每个root处res加上当前值异或父节点的值"""
            if root >= len(preOrder):
                return

            rootVal = preOrder[self.index]
            self.index += 1
            valueByRoot[root] = rootVal
            if root >= 1:
                self.res += rootVal ^ valueByRoot[(root - 1) >> 1]

            dfs((root << 1) + 1)
            dfs((root << 1) + 2)

        # write code here
        valueByRoot = defaultdict(int)
        self.res = 0
        self.index = 0
        dfs(0)
        return self.res


print(1)
