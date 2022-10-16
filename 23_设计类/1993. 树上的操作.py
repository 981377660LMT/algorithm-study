from typing import List
from collections import defaultdict, deque

# 给你一棵 n 个节点的树，编号从 0 到 n - 1 ，以父节点数组 parent 的形式给出，其中 parent[i] 是第 i 个节点的父节点


class LockingTree:
    def __init__(self, parent: List[int]):
        self.parent = parent
        self.adjMap = defaultdict(set)
        for i, pre in enumerate(parent):
            if pre != -1:
                self.adjMap[pre].add(i)
        self.locked = dict()

    # 指定用户给指定节点 上锁 ，上锁后其他用户将无法给同一节点上锁。只有当节点处于未上锁的状态下，才能进行上锁操作。
    def lock(self, num: int, user: int) -> bool:
        if num in self.locked:
            return False
        self.locked[num] = user
        return True

    # 指定用户给指定节点 解锁 ，只有当指定节点当前正被指定用户锁住时，才能执行该解锁操作。
    def unlock(self, num: int, user: int) -> bool:
        if self.locked.get(num, -1) != user:
            return False
        self.locked.pop(num)
        return True

    # !指定用户给指定节点 上锁 ，并且将该节点的所有子孙节点 解锁 。
    # 升级条件：
    # 指定节点当前状态为未上锁。
    # 指定节点至少有一个上锁状态的子孙节点（可以是 任意 用户上锁的）。
    # 指定节点没有任何上锁的祖先节点。
    def upgrade(self, num: int, user: int) -> bool:
        if num in self.locked:
            return False

        # 检查没有任何上锁的祖先结点
        root = num
        while root != -1:
            if root in self.locked:
                return False
            root = self.parent[root]

        # 检查至少有一个上锁的子孙结点
        queue = deque([num])
        lockedChildren = []
        while queue:
            cur = queue.popleft()
            if cur in self.locked:
                lockedChildren.append(cur)
            for child in self.adjMap[cur]:
                queue.append(child)

        if not lockedChildren:
            return False

        # 要求的操作
        self.locked[num] = user
        for child in lockedChildren:
            self.locked.pop(child)
        return True
