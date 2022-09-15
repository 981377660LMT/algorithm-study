# https://leetcode-cn.com/problems/longest-substring-of-one-repeating-character/solution/python-guo-ran-wo-huan-shi-geng-xi-huan-olhop/


from typing import List


# 1 <= s.length <= 1e5
# 1 <= k <= 1e5

# 两处优化
# 1. build 建树后不需要update
# 2. Node 换成数组


class SegmentTree:
    def __init__(self, s: str) -> None:
        n = len(s)
        self.tree = [0 for _ in range(n * 4)]  # [left,right]区间内的最大连续数
        self.pre = [0 for _ in range(n * 4)]  # 区间左端点为起点的连续数
        self.suf = [0 for _ in range(n * 4)]  # 区间右端点为终点的连续数
        self.chars = list(s)
        self.build(1, 1, n)

    def build(self, rt: int, left: int, right: int) -> None:
        if left == right:
            self.tree[rt] = 1
            self.pre[rt] = 1
            self.suf[rt] = 1
            return

        mid = (left + right) // 2
        self.build(rt * 2, left, mid)
        self.build(rt * 2 + 1, mid + 1, right)
        self._pushUp(rt, left, right, mid)

    def update(self, rt: int, L: int, R: int, left: int, right: int, target: str) -> None:
        """区间修改，L,R表示需要update的范围,l,r表示当前节点的范围"""
        """===> 如果这里是真正的区间修改应该带有lazy标记，tree pre suf什么的都等于right-left+1 <==="""
        if L <= left and right <= R:
            self.tree[rt] = 1
            self.pre[rt] = 1
            self.suf[rt] = 1
            self.chars[left - 1] = target
            return

        mid = (left + right) // 2
        if L <= mid:
            self.update(rt * 2, L, R, left, mid, target)
        if mid + 1 <= R:
            self.update(rt * 2 + 1, L, R, mid + 1, right, target)
        self._pushUp(rt, left, right, mid)

    def query(self, rt: int, L: int, R: int, left: int, right: int) -> int:
        """L,R表示需要query的范围,left,right表示当前节点的范围"""
        """===> 虽然但是，如果要写真正的区间查询，应该按照类似pushup的方法把一堆结果进行合并 <==="""
        if L <= left and right <= R:
            return self.tree[rt]

        mid = (left + right) // 2
        res = 0
        if L <= mid:
            res = max(res, self.query(rt * 2, L, R, left, mid))
        if mid + 1 <= R:
            res = max(res, self.query(rt * 2 + 1, L, R, mid + 1, right))
        return res

    def queryAll(self) -> int:
        return self.tree[1]

    def _pushUp(self, rt: int, left: int, right: int, mid: int) -> None:
        leftPre, rightPre = self.pre[rt * 2], self.pre[rt * 2 + 1]
        leftSuf, rightSuf = self.suf[rt * 2], self.suf[rt * 2 + 1]
        leftMax, rightMax = self.tree[rt * 2], self.tree[rt * 2 + 1]

        self.pre[rt] = leftPre
        self.suf[rt] = rightSuf

        if self.chars[mid - 1] == self.chars[mid]:
            # 合并
            self.tree[rt] = max(leftMax, rightMax, leftSuf + rightPre)
            if leftPre == mid - left + 1:
                self.pre[rt] += rightPre
            if rightSuf == right - mid:
                self.suf[rt] += leftSuf
        else:
            self.tree[rt] = max(leftMax, rightMax)


class Solution:
    def longestRepeating(self, s: str, queryCharacters: str, queryIndices: List[int]) -> List[int]:
        """
        第 i 个查询会将 s 中位于下标 queryIndices[i] 的字符更新为 queryCharacters[i] 。
        返回一个长度为 k 的数组 lengths ，其中 lengths[i] 是在执行第 i 个查询 之后 s 中仅由 单个字符重复 组成的 最长子字符串 的 长度 。
        """
        n = len(s)
        segmentTree = SegmentTree(s)
        res = [0] * len(queryIndices)

        for index, (qc, qi) in enumerate(zip(queryCharacters, queryIndices)):
            segmentTree.update(1, qi + 1, qi + 1, 1, n, qc)
            # 因为每次query整个线段树区间，所以不要懒更新
            res[index] = segmentTree.queryAll()
        return res


print(Solution().longestRepeating(s="babacc", queryCharacters="bcb", queryIndices=[1, 3, 3]))
# [3,3,4]
