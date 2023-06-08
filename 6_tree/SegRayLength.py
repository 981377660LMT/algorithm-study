# https://github.com/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/TREE/SegRayLengthHelper_vector.md#L3
# 最长线段射线助手。

# seg:线段,这里指的是某条路径.
# ray:射线,这里指的是以某个树节点为起点的路径.
# 对于根结点，所有的连通部分都在自己下方。
# 对于非根结点，一个连通部分在自己上方，剩余连通部分都在自己下方。
# 在预处理后，
# !  `ray[i]` 存储了 以 `i` 作为一端的线段的最长长度`前三名`。
#   `seg[i]` 存储了不经过 `i` 的线段的最长长度前两名。
#   `downRay[i]` 存储了以 `i` 作为一端的向下的线段的最长长度。
#   `downSeg[i]` 存储了 `i` 下方的不越过 `i` 的线段的线段最长长度。注意，不越过，但是可以经过。
#   `upRay[i]` 存储了 以 `i` 作为一端的向上的线段的最长长度。
#   `upSeg[i]` 存储了 `i` 上方的不越过 `i` 的线段的线段最长长度。注意，不越过，但是可以经过。
#
# !注意：此处的前三名、前两名不允许某个连通部分(子树)重复占用。
# 比如结点 `i` 有两个邻居，从第一个连通部分可以找到长度为 `10` 的 `ray` ，也可以找到长度为 `9` 的 `ray` ；
# 从第二个联通部分可以找到长度为 `8` 的 `ray` ，
# 那么 `m_ray[i][0]==10` ，`m_ray[i][1]==8` ，`m_ray[i][2]==0` ，也就是没有第三名。


from typing import List, Tuple


class SegRayLength:
    """
    用于处理树中最长射线/线段的问题.
    https://github.com/old-yan/CP-template/blob/main/TREE/SegRayLengthHelper_vector.h
    """

    __slots__ = (
        "ray",
        "seg",
        "downRay",
        "downSeg",
        "upRay",
        "upSeg",
        "_tree",
        "_depthWeighted",
        "_lid",
        "_top",
        "_parent",
        "_heavySon",
        "_dfn",
    )

    def __init__(self, n: int) -> None:
        self.ray = [[0, 0, 0] for _ in range(n)]
        self.seg = [[0, 0] for _ in range(n)]
        self.downRay = [0 for _ in range(n)]
        self.downSeg = [0 for _ in range(n)]
        self.upRay = [0 for _ in range(n)]
        self.upSeg = [0 for _ in range(n)]
        self._tree = [[] for _ in range(n)]  # (next,weight)

        # HLD
        self._depthWeighted = [0] * n
        self._lid = [0] * n
        self._top = [0] * n
        self._parent = [-1] * n
        self._heavySon = [0] * n
        self._dfn = 0

    def addEdge(self, u: int, v: int, w=1) -> None:
        self.addDirectedEdge(u, v, w)
        self.addDirectedEdge(v, u, w)

    def addDirectedEdge(self, u: int, v: int, w=1) -> None:
        self._tree[u].append((v, w))

    def build(self, root=0) -> None:
        self._dfs1(root, -1, 0)
        self._dfs2(root, -1, 0, 0)
        self._dfs3(root, root)  # markTop

    def queryMaxRayAndSeg(self, u: int, ignoreRoot=-1) -> Tuple[int, int]:
        """
        查询最长射线ray,最长线段seg.

        Args:
            u:树节点.
            ignoreRoot:屏蔽掉以ignoreRoot为根的子树.需要保证ignoreRoot在u的子树中.

        Returns:
            Tuple[int,int]:树中剩余部分从u出发的最长射线,树中剩余部分的最长线段.
        """
        if ignoreRoot == -1:
            return self._maxRaySeg(u, -1, -1)
        return self._maxRaySeg(
            u,
            self.downRay[ignoreRoot] + self._weightedDist(u, ignoreRoot),
            self.downSeg[ignoreRoot],
        )

    def _dfs1(self, cur: int, pre: int, dist: int) -> int:
        subSize, heavySize, heavySon = 1, 0, -1
        downRay, downSeg = self.downRay, self.downSeg
        for next, weight in self._tree[cur]:
            if next != pre:
                nextSize = self._dfs1(next, cur, dist + weight)
                subSize += nextSize
                if nextSize > heavySize:
                    heavySize, heavySon = nextSize, next
                self._addDownRay(cur, downRay[next] + weight)
                self._addDownSeg(cur, downSeg[next])
        len1, len2, _ = self.ray[cur]
        cand = len1 + len2
        if cand > downSeg[cur]:
            downSeg[cur] = cand
        self._depthWeighted[cur] = dist
        self._parent[cur] = pre
        self._heavySon[cur] = heavySon
        return subSize

    def _dfs2(self, cur: int, pre: int, upRay: int, upSeg: int) -> None:
        self._setUpRay(cur, upRay)
        self._setUpSeg(cur, upSeg)
        downRay, downSeg = self.downRay, self.downSeg
        for next, weight in self._tree[cur]:
            if next != pre:
                ray, seg = self._maxRaySeg(cur, downRay[next] + weight, downSeg[next])
                self._addSeg(next, seg)
                ray += weight
                if ray > seg:
                    seg = ray
                self._dfs2(next, cur, ray, seg)

    def _dfs3(self, cur: int, top: int) -> None:
        self._top[cur] = top
        self._lid[cur] = self._dfn
        self._dfn += 1
        heavySon = self._heavySon[cur]
        if heavySon != -1:
            self._dfs3(heavySon, top)
            for next, _ in self._tree[cur]:
                if next != heavySon and next != self._parent[cur]:
                    self._dfs3(next, next)

    def _addRay(self, i: int, ray: int) -> None:
        rayI = self.ray[i]
        if ray > rayI[0]:
            rayI[2] = rayI[1]
            rayI[1] = rayI[0]
            rayI[0] = ray
        elif ray > rayI[1]:
            rayI[2] = rayI[1]
            rayI[1] = ray
        elif ray > rayI[2]:
            rayI[2] = ray

    def _addSeg(self, i: int, seg: int) -> None:
        segI = self.seg[i]
        if seg > segI[0]:
            segI[1] = segI[0]
            segI[0] = seg
        elif seg > segI[1]:
            segI[1] = seg

    def _addDownRay(self, i: int, ray: int) -> None:
        if ray > self.downRay[i]:
            self.downRay[i] = ray
        self._addRay(i, ray)

    def _addDownSeg(self, i: int, seg: int) -> None:
        if seg > self.downSeg[i]:
            self.downSeg[i] = seg
        self._addSeg(i, seg)

    def _setUpRay(self, i: int, ray: int) -> None:
        self.upRay[i] = ray
        self._addRay(i, ray)

    def _setUpSeg(self, i: int, seg: int) -> None:
        self.upSeg[i] = seg

    def _maxRaySeg(self, u: int, ignoreRay: int, ignorSeg: int) -> Tuple[int, int]:
        """
        查询树中某部分的最长射线和线段.

        Args:
            u:树中某点.
            ignoreRay:屏蔽掉的部分提供的最长射线.
            ignorSeg:屏蔽掉的部分提供的最长线段.

        Returns:
            Tuple[int, int]:从u出发的最长射线和树中剩余部分的最长线段.
        """
        r0, r1, r2 = self.ray[u]
        s0, s1 = self.seg[u]
        maxRay = r1 if ignoreRay == r0 else r0
        maxSeg = s1 if ignorSeg == s0 else s0
        twoRay = r1 + r2 if ignoreRay == r0 else (r0 + r2 if ignoreRay == r1 else r0 + r1)
        if maxSeg < twoRay:
            maxSeg = twoRay
        return maxRay, maxSeg

    def _lca(self, u: int, v: int) -> int:
        while True:
            if self._lid[u] > self._lid[v]:
                u, v = v, u
            if self._top[u] == self._top[v]:
                return u
            v = self._parent[self._top[v]]

    def _weightedDist(self, u: int, v: int) -> int:
        return (
            self._depthWeighted[u]
            + self._depthWeighted[v]
            - 2 * self._depthWeighted[self._lca(u, v)]
        )


if __name__ == "__main__":
    n = 5
    SR = SegRayLength(n)
    SR.addEdge(2, 0, 1)
    SR.addEdge(1, 3, 1)
    SR.addEdge(4, 0, 1)
    SR.addEdge(3, 0, 1)
    SR.build(root=3)

    print(SR.ray)
    print(SR.seg)
    print(SR.downRay)
    print(SR.downSeg)
    print(SR.upRay)
    print(SR.upSeg)

    assert SR.queryMaxRayAndSeg(0, 2) == (2, 3)

    class Solution:
        def treeDiameter(self, edges: List[List[int]]) -> int:
            SR = SegRayLength(len(edges) + 1)
            for u, v in edges:
                SR.addEdge(u, v, 1)
            SR.build(root=0)
            return max(SR.ray[i][0] for i in range(len(edges) + 1))
