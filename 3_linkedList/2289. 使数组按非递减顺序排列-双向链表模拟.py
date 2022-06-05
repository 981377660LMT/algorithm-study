# 给你一个下标从 0 开始的整数数组 nums 。在一步操作中，移除所有满足 nums[i - 1] > nums[i] 的 nums[i] ，其中 0 < i < nums.length 。
# 重复执行步骤，直到 nums 变为 非递减 数组，返回所需执行的操作数。


# 每一轮操作中，我们按 下标 i 递减的顺序 遍历列表 v，对于其中的每个下标 i，
# 如果 nums[i] 没有被吃掉，且 nums[i] > 它的 后继元素，
# 那么 删除 nums[i] 的 后继元素。
# 另外，nums[i] 在下一轮的操作中可能继续吃掉其它数字，将其保存到一个新的列表 v2 中。
# 否则，要么 nums[i] 被吃掉，要么 nums[i]≤ 它的后继元素。在
# 这两种情况下，nums[i] 都再也不能吃掉掉后面的数字了，因此直接丢弃即可。

# !链表模拟

from collections import deque
from dataclasses import dataclass
from typing import List, Optional


@dataclass(slots=True)
class Node:
    value: int
    left: Optional['Node']
    right: Optional['Node']
    hasDel: bool = False

    def __repr__(self) -> str:
        return f'{self.value}'


class Solution:
    def totalSteps(self, nums: List[int]) -> int:
        def remove(node: Optional['Node']) -> None:
            if not node:
                return
            if node.left:
                node.left.right = node.right
            if node.right:
                node.right.left = node.left
            node.hasDel = True

        # 尾部加入一个哨兵
        # nums.append(int(1e20))
        n = len(nums)
        nodes = [Node(value, None, None) for value in nums]
        for i in range(n):
            if i + 1 < n:
                nodes[i].right = nodes[i + 1]
            if i - 1 >= 0:
                nodes[i].left = nodes[i - 1]

        queue = deque([node for node in nodes if (node.right and node.value > node.right.value)])
        res = 0
        while queue:
            # print(queue)
            hasRemoved = False
            nextQueue = deque()
            for node in reversed(queue):
                if node.hasDel:
                    continue
                if node.right and node.value > node.right.value:
                    remove(node.right)
                    nextQueue.appendleft(node)
                    hasRemoved = True
            queue = nextQueue
            res += 1 if hasRemoved else 0

        return res


print(Solution().totalSteps(nums=[5, 3, 4, 4, 7, 3, 6, 11, 8, 5, 11]))
