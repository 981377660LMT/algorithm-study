# 问题就简化为：求在n个首尾相连的元素中，选取n/3个不相邻元素所能获得的最大值。
# 一共有n+n/3次入队操作，所以总的时间复杂度为O(nlogn)。

# 每次选择了一个节点，然后调整该节点的权值，相当于将原问题转化为规模为 N-2 的子问题
# 比如 [1,1,1,8,9,8,6,1,1] 先选 9，此时两边的 8 不能选。
# 但若选了两边的 8，则 ‘8 9 8’ 两边的 ‘1’ ‘6’ 又不能选。如何简化这个逻辑？
# 如果这样，选择 9 后，将 9 调整为 8 + 8 - 9 = 7，则只需解决子问题 [1,1,1,'‘7'’,6,1,1] 即可：
# 若选 7，相当于选了原数组的两个 ‘8'，同时也造成了 '7' 两边的 '1' '6' 不能选，
# 并且批萨的数量也增加了 1。这样就完美简化了上面的逻辑。
# !1388. 3n 块披萨-优先队列+双向链表

from dataclasses import dataclass
from heapq import heapify, heappop, heappush
from typing import List


@dataclass(slots=True)
class MaxCycleNode:
    value: int
    left: 'MaxCycleNode'
    right: 'MaxCycleNode'
    hasDel: bool = False

    def __eq__(self, other: 'MaxCycleNode') -> bool:
        return self.value == other.value

    def __lt__(self, other: 'MaxCycleNode') -> bool:
        return self.value > other.value

    def __repr__(self) -> str:
        return f'{self.value}'


class Solution:
    def maxSizeSlices(self, slices: List[int]) -> int:
        def remove(node: 'MaxCycleNode') -> None:
            node.left.right = node.right
            node.right.left = node.left
            node.hasDel = True  # 标记删除

        n = len(slices)
        pq = [MaxCycleNode(value, None, None, False) for value in slices]  # type: ignore
        for i in range(n):
            pq[i].left = pq[(i - 1) % n]
            pq[i].right = pq[(i + 1) % n]
        heapify(pq)

        res = []
        while len(res) < n // 3:
            maxNode = heappop(pq)
            if maxNode.hasDel:  # !堆的标记删除
                continue
            res.append(maxNode.value)
            maxNode.value = maxNode.left.value + maxNode.right.value - maxNode.value
            remove(maxNode.left)
            remove(maxNode.right)
            heappush(pq, maxNode)

        return sum(res)


print(Solution().maxSizeSlices(slices=[1, 2, 3, 4, 5, 6]))
