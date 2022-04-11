from typing import List, Tuple

# 子树结点颜色都不同的结点个数
class Solution:
    def solve(self, tree: List[List[int]], color: List[int]) -> int:
        """位运算压缩颜色"""

        def dfs(cur: int, parent: int) -> Tuple[bool, int]:
            nonlocal res
            isSpecial = True
            subTree = 1 << color[cur]
            for next in tree[cur]:
                if next == parent:
                    continue

                nextSpecial, nextSubTree = dfs(next, cur)
                if not nextSpecial or nextSubTree & subTree:
                    isSpecial = False
                subTree |= nextSubTree

            if isSpecial:
                res += 1

            return isSpecial, subTree

        res = 0
        dfs(0, -1)
        return res

