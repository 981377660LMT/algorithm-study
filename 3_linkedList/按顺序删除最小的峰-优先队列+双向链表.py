# 峰是比两边都大的数
# 按顺序删除最小的峰


from heapq import heapify, heappop, heappush
from typing import List, Optional


class MinCycleNode:
    __slots__ = ("index", "value", "left", "right", "deleted")

    def __init__(
        self,
        index: int,
        value: int,
        left: Optional["MinCycleNode"] = None,
        right: Optional["MinCycleNode"] = None,
    ) -> None:
        self.index = index
        self.value = value
        self.left = left
        self.right = right
        self.deleted = False

    def __eq__(self, other: "MinCycleNode") -> bool:
        return self.value == other.value

    def __lt__(self, other: "MinCycleNode") -> bool:
        return self.value < other.value

    def __repr__(self) -> str:
        return f"{self.index} {self.value} {self.deleted}"


def remove(node: Optional["MinCycleNode"]) -> None:
    if node is None:
        return
    if node.left:
        node.left.right = node.right
    if node.right:
        node.right.left = node.left
    node.deleted = True  # 标记删除


# !优先队列+循环双向链表
class Solution:
    def solve(self, nums: List[int]) -> List[int]:
        nums.append(-int(1e20))  # 哨兵
        n = len(nums)
        nodes = [MinCycleNode(index, value) for index, value in enumerate(nums)]
        for i in range(n):  # !双向循环链表
            nodes[i].left = nodes[(i - 1) % n]
            nodes[i].right = nodes[(i + 1) % n]
        peaks = [node for node in nodes if node.value > max(node.left.value, node.right.value)]  # type: ignore
        heapify(peaks)

        res = []
        while len(res) < n - 1:
            minPeak = heappop(peaks)
            res.append(minPeak.value)
            remove(minPeak)
            # !删除一个峰后 影响到了相邻的峰候选 要用链表获取前驱后继
            for cand in (minPeak.left, minPeak.right):
                if cand.value > max(cand.left.value, cand.right.value):
                    heappush(peaks, cand)
        return res


print(Solution().solve(nums=[3, 5, 1, 4, 2]))

# [4, 2, 5, 3, 1]
# We remove 4 and get [3, 5, 1, 2]
# We remove 2 and get [3, 5, 1]
# We remove 5 and get [3, 1]
# We remove 3 and get [1]
# We remove 1 and get []
