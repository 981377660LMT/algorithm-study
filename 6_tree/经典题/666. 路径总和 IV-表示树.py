from typing import List, Tuple

# 对于一棵深度小于 5 的树，可以用一组三位十进制整数来表示。
# 深度，层所在位置，权值
# 给定一个包含三位整数的升序数组，表示一棵深度小于 5 的二叉树，请你返回从根到所有叶子结点的路径之和。

# 总结：
# 1.如果root为(10*depth,pos) 那么左节点为(10*(depth+1),2*pos-1) 右节点为(10*(depth+1),2*pos)
# 2.使用{id:val} 来存储结点 使用divmod(root) 来获取结点层数与depth


class Solution:
    def pathSum(self, nums: List[int]) -> int:
        self.res = 0
        valueByNode = {n // 10: n % 10 for n in nums}  # 结点id:结点的值

        def dfs(root: int, s: int):
            if root not in valueByNode:
                return

            s += valueByNode[root]
            depth, pos = divmod(root, 10)
            left = (depth + 1) * 10 + 2 * pos - 1
            right = left + 1

            if left not in valueByNode and right not in valueByNode:
                self.res += s
            else:
                dfs(left, s)
                dfs(right, s)

        dfs(nums[0] // 10, 0)
        return self.res


print(Solution().pathSum([113, 215, 221]))
# 输出: 12
# 解释:
# 这棵树形状如下:
#     3
#    / \
#   5   1
# 路径和 = (3 + 5) + (3 + 1) = 12.
