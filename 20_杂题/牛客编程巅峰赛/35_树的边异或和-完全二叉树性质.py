from collections import defaultdict

# 完全二叉树


class Solution:
    def tree5(self, preOrder: list[int]) -> int:
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
