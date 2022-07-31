from typing import List

INF = int(1e20)


class Node:
    __slots__ = ("left", "right", "sum", "max", "lMax", "rMax")

    def __init__(self) -> None:
        self.left = -1
        self.right = -1
        self.sum = 0  # 区间和
        self.max = -INF  # 区间最大子数组和
        self.lMax = -INF  # 区间左端点的最大子数组和
        self.rMax = -INF  # 区间右端点的最大子数组和


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    def __init__(self, nums: List[int]):
        n = len(nums)
        self.tree = [Node() for _ in range(n << 2)]
        self.build(1, 1, n, nums)

    def build(self, rt: int, left: int, right: int, nums: List[int]) -> None:
        root = self.tree[rt]
        root.left = left
        root.right = right
        if left == right:
            root.sum = nums[left - 1]
            root.max = nums[left - 1]
            root.lMax = nums[left - 1]
            root.rMax = nums[left - 1]
            return

        mid = (left + right) >> 1
        self.build(rt << 1, left, mid, nums)
        self.build(rt << 1 | 1, mid + 1, right, nums)
        self.pushUp(rt)

    def query(self, rt: int, left: int, right: int) -> int:
        """区间和"""
        root = self.tree[rt]
        if left <= root.left and root.right <= right:
            return root.max

        res = -INF
        mid = (root.left + root.right) >> 1
        if left <= mid:
            res = max(res, self.query(rt << 1, left, right))
        if mid < right:
            res = max(res, self.query(rt << 1 | 1, left, right))
        return res

    def update(self, rt: int, left: int, right: int, val: int) -> None:
        """单点更新"""
        root = self.tree[rt]
        if left <= root.left and root.right <= right:
            root.sum = val
            root.max = val
            root.lMax = val
            root.rMax = val
            return

        mid = (root.left + root.right) >> 1
        if left <= mid:
            self.update(rt << 1, left, right, val)
        if mid < right:
            self.update(rt << 1 | 1, left, right, val)
        self.pushUp(rt)

    def pushUp(self, rt: int) -> None:
        """不要懒更新"""
        root, left, right = self.tree[rt], self.tree[rt << 1], self.tree[rt << 1 | 1]
        root.sum = left.sum + right.sum
        root.lMax = max(left.lMax, left.sum + right.lMax)
        root.rMax = max(right.max, right.sum + left.rMax)
        root.max = max(left.rMax + right.lMax, left.max, right.max)


n, m = map(int, input().split())
nums = list(map(int, input().split()))
tree = SegmentTree(nums)
res = []

for _ in range(m):
    opt, *rest = input().split()
    if opt == "1":
        l, r = sorted(map(int, rest))
        # 查询区间 [x,y] 中的最大子数组和
        res.append(str(tree.query(1, l, r)))
    else:
        l, val = map(int, rest)
        # 把 A[x] 改成 y。
        tree.update(1, l, l, val)


print("\n".join(res))
