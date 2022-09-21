"""
某次篮球比赛前,
需要将站成一排的n名球员分为两队(两队人数可以不同),
每名球员的能力值为ai。
有两名教练轮流挑选队员,第一个教练先挑选。
每位教练每次选人时,都会选择当前剩余的所有人中,
能力值最大的那一个。当选择一个人后,
会将他左右两侧各m个人一起挑选走(若某一侧可选的人数不够m人,则将这—侧能选的人都选上)。
请输出此规则下,分到两队的具体成员情况。
"""

from heapq import heapify, heappop
from typing import List, Optional


class MaxCycleNode:
    __slots__ = ("index", "value", "left", "right", "deleted")

    def __init__(
        self,
        index: int,
        value: int,
        left: Optional["MaxCycleNode"] = None,
        right: Optional["MaxCycleNode"] = None,
    ) -> None:
        self.index = index
        self.value = value
        self.left = left
        self.right = right
        self.deleted = False

    def __eq__(self, other: "MaxCycleNode") -> bool:
        return self.value == other.value

    def __lt__(self, other: "MaxCycleNode") -> bool:
        return self.value > other.value

    def __repr__(self) -> str:
        return f"{self.index} {self.value} {self.deleted}"


def remove(node: Optional["MaxCycleNode"]) -> None:
    if node is None:
        return
    if node.left:
        node.left.right = node.right
    if node.right:
        node.right.left = node.left
    node.deleted = True  # 标记删除


def solve(n: int, m: int, nums: List[int]) -> List[str]:
    def select(team: str) -> None:
        maxNode = None
        while pq:
            cur = heappop(pq)
            if not cur.deleted:
                res[cur.index] = team
                maxNode = cur
                break

        if maxNode is None:
            return

        left, right = maxNode.left, maxNode.right
        remove(maxNode)

        count = m
        while count > 0 and left:
            res[left.index] = team
            remove(left)
            left = left.left
            count -= 1

        count = m
        while count > 0 and right:
            res[right.index] = team
            remove(right)
            right = right.right
            count -= 1

    res = [""] * n
    pq = [MaxCycleNode(index, value) for index, value in enumerate(nums)]
    for i in range(n):  # 双向链表
        if i - 1 >= 0:
            pq[i].left = pq[(i - 1)]
        if i + 1 < n:
            pq[i].right = pq[(i + 1)]

    heapify(pq)

    while pq:
        select("A")
        select("B")

    return res


if __name__ == "__main__":
    assert solve(7, 1, [3, 6, 1, 7, 2, 5, 4]) == ["B", "B", "A", "A", "A", "B", "A"]
    assert solve(10, 2, [4, 8, 9, 10, 7, 6, 5, 3, 2, 1]) == [
        "B",
        "A",
        "A",
        "A",
        "A",
        "A",
        "B",
        "B",
        "B",
        "A",
    ]
