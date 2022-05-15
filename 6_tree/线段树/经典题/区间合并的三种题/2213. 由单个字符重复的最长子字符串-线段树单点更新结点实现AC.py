# https://leetcode-cn.com/problems/longest-substring-of-one-repeating-character/solution/python-guo-ran-wo-huan-shi-geng-xi-huan-olhop/


from typing import List


# 1 <= s.length <= 1e5
# 1 <= k <= 1e5
class Node:
    __slots__ = ('left', 'right', 'pre', 'suf', 'max')

    def __init__(self, left=-1, right=-1) -> None:
        self.left = left
        self.right = right
        self.pre = 1  # 区间左端点的连续数
        self.suf = 1  # 区间右端点的连续数
        self.max = 1  # [left,right]区间内的最大连续数


class SegmentTree:
    def __init__(self, s: str) -> None:
        n = len(s)
        self.tree = [Node() for _ in range(n << 2)]
        self.chars = list(s)
        self.build(1, 1, n)

    def build(self, rt: int, left: int, right: int) -> None:
        root = self.tree[rt]
        root.left, root.right = left, right
        if left == right:
            root.pre = 1
            root.suf = 1
            root.max = 1
            return

        mid = (root.left + root.right) >> 1
        self.build(rt << 1, left, mid)
        self.build(rt << 1 | 1, mid + 1, right)
        self.pushUp(rt)

    def update(self, rt: int, left: int, right: int, target: str) -> None:
        """区间修改，L,R表示需要update的范围,l,r表示当前节点的范围"""
        """===> 如果这里是真正的区间修改应该带有lazy标记，tree pre suf什么的都等于right-left+1 <==="""
        root = self.tree[rt]
        if left <= root.left <= root.right <= right:
            root.pre = 1
            root.suf = 1
            root.max = 1
            # 注意这里修改
            self.chars[left - 1] = target
            return

        mid = (root.left + root.right) >> 1
        if left <= mid:
            self.update(rt << 1, left, right, target)
        if mid < right:
            self.update(rt << 1 | 1, left, right, target)
        self.pushUp(rt)

    def query(self, rt: int, left: int, right: int) -> int:
        """L,R表示需要query的范围,left,right表示当前节点的范围"""
        """===> 虽然但是，如果要写真正的区间查询，应该按照类似pushup的方法把一堆结果进行合并 <==="""

        root = self.tree[rt]
        if left <= root.left and root.right <= right:
            return root.max

        mid = (root.left + root.right) >> 1
        res = 0
        if left <= mid:
            res = max(res, self.query(rt << 1, left, right))
        if mid < right:
            res = max(res, self.query(rt << 1 | 1, left, right))
        return res

    def pushUp(self, rt: int) -> None:
        root, left, right = self.tree[rt], self.tree[(rt << 1)], self.tree[(rt << 1) | 1]
        root.pre = left.pre
        root.suf = right.suf

        mid = (root.left + root.right) >> 1
        if self.chars[mid - 1] == self.chars[mid]:
            # 合并
            root.max = max(left.max, right.max, left.suf + right.pre)
            if left.pre == left.right - left.left + 1:
                root.pre += right.pre
            if right.suf == right.right - right.left + 1:
                root.suf += left.suf
        else:
            root.max = max(left.max, right.max)


class Solution:
    def longestRepeating(self, s: str, queryCharacters: str, queryIndices: List[int]) -> List[int]:
        """
        第 i 个查询会将 s 中位于下标 queryIndices[i] 的字符更新为 queryCharacters[i] 。
        返回一个长度为 k 的数组 lengths ，其中 lengths[i] 是在执行第 i 个查询 之后 s 中仅由 单个字符重复 组成的 最长子字符串 的 长度 。
        """
        segmentTree = SegmentTree(s)
        res = [0] * len(queryIndices)
        for index, (qc, qi) in enumerate(zip(queryCharacters, queryIndices)):
            segmentTree.update(1, qi + 1, qi + 1, qc)
            # 因为每次query整个线段树区间，所以不要懒更新
            res[index] = segmentTree.query(1, 1, len(s))
        return res


print(Solution().longestRepeating(s="babacc", queryCharacters="bcb", queryIndices=[1, 3, 3]))
# [3,3,4]

