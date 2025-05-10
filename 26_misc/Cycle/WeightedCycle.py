from itertools import accumulate, chain, pairwise
from typing import List, Tuple, Optional
from bisect import bisect_left, bisect_right


def min2(a: int, b: int) -> int:
    return a if a < b else b


class WeightedCycle:
    """带权重的环类，使用前缀和优化查询，支持有向边权重."""

    __slots__ = "n", "weights_cw", "weights_ccw", "prefix_cw", "prefix_ccw"

    def __init__(
        self,
        n: int,
        weights_cw: Optional[List[int]] = None,
        weights_ccw: Optional[List[int]] = None,
    ):
        """
        - 默认所有边权重为1

        :param n: 环上的节点数量

        :param weights_cw: 顺时针方向边的权重列表，如果提供，长度应为n
                          weights_cw[i]表示从i顺时针到(i+1)%n的边权重

        :param weights_ccw: 逆时针方向边的权重列表，如果提供，长度应为n
                          weights_ccw[i]表示从(i+1)%n逆时针到i的边权重
        """
        self.n = n

        if weights_cw is None:
            self.weights_cw = [1] * n
        elif len(weights_cw) == n:
            self.weights_cw = weights_cw
        else:
            raise ValueError(f"顺时针权重数量应为{n}，但提供了{len(weights_cw)}个")

        if weights_ccw is None:
            self.weights_ccw = [1] * n
        elif len(weights_ccw) == n:
            self.weights_ccw = weights_ccw
        else:
            raise ValueError(f"逆时针权重数量应为{n}，但提供了{len(weights_ccw)}个")

        self._build_prefix_sums()

    def _build_prefix_sums(self):
        self.prefix_cw = list(accumulate(self.weights_cw, initial=0))
        self.prefix_ccw = list(accumulate(self.weights_ccw, initial=0))

    def dist(self, u: int, v: int) -> int:
        return min2(self.dist_ccw(u, v), self.dist_cw(u, v))

    def dist_ccw(self, from_: int, to: int) -> int:
        """逆时针从from_到to的带权距离."""
        if from_ >= to:
            return self.prefix_ccw[from_] - self.prefix_ccw[to]
        else:
            return self.prefix_ccw[from_] + (self.prefix_ccw[self.n] - self.prefix_ccw[to])

    def dist_cw(self, from_: int, to: int) -> int:
        """顺时针从from_到to的带权距离."""
        if to >= from_:
            return self.prefix_cw[to] - self.prefix_cw[from_]
        else:
            return (self.prefix_cw[self.n] - self.prefix_cw[from_]) + self.prefix_cw[to]

    def segment(self, u: int, v: int) -> List[Tuple[int, int]]:
        """返回环上两点间最短路径的线段表示"""
        if self.dist_ccw(u, v) <= self.dist_cw(u, v):
            return self.segment_ccw(u, v)
        return self.segment_cw(u, v)

    def segment_ccw(self, from_: int, to: int) -> List[Tuple[int, int]]:
        """逆时针从from_到to的路径段"""
        if from_ >= to:
            return [(from_, to)]
        return [(from_, 0), (self.n - 1, to)]

    def segment_cw(self, from_: int, to: int) -> List[Tuple[int, int]]:
        """顺时针从from_到to的路径段"""
        if to >= from_:
            return [(from_, to)]
        return [(from_, self.n - 1), (0, to)]

    def path(self, u: int, v: int) -> List[int]:
        """返回环上两点间的最短路径"""
        if self.dist_ccw(u, v) <= self.dist_cw(u, v):
            return self.path_ccw(u, v)
        return self.path_cw(u, v)

    def path_ccw(self, from_: int, to: int) -> List[int]:
        """逆时针从from_到to的路径经过的点"""
        if from_ >= to:
            return list(range(from_, to - 1, -1))
        return list(range(from_, -1, -1)) + list(range(self.n - 1, to - 1, -1))

    def path_cw(self, from_: int, to: int) -> List[int]:
        """顺时针从from_到to的路径经过的点"""
        if to >= from_:
            return list(range(from_, to + 1))
        return list(range(from_, self.n)) + list(range(0, to + 1))

    def on_path_ccw(self, from_: int, to: int, x: int) -> bool:
        """x是否在from_到to的逆时针路径上"""
        if x < to:
            x += self.n
        if from_ < to:
            from_ += self.n
        return to <= x <= from_

    def on_path_cw(self, from_: int, to: int, x: int) -> bool:
        """x是否在from_到to的顺时针路径上"""
        if from_ > to:
            to += self.n
        if from_ > x:
            x += self.n
        return from_ <= x <= to

    def jump_ccw(self, from_: int, distance: int) -> int:
        """
        逆时针从from_出发走特定距离到达的位置.
        """
        if distance == 0:
            return from_

        total_weight = self.prefix_ccw[self.n]
        distance %= total_weight
        if distance == 0:
            return from_

        target = self.prefix_ccw[from_] - distance
        if target < 0:
            target += total_weight

        pos = bisect_left(self.prefix_ccw, target)
        return pos if pos < self.n else 0

    def jump_cw(self, from_: int, distance: int) -> int:
        """
        顺时针从from_出发走特定距离到达的位置.
        只有当距离完全达到时才能抵达下一个点.
        """
        if distance == 0:
            return from_

        total_weight = self.prefix_cw[self.n]
        distance %= total_weight
        if distance == 0:
            return from_

        target = self.prefix_cw[from_] + distance
        if target >= total_weight:
            target -= total_weight

        pos = bisect_right(self.prefix_cw, target)
        return pos - 1


if __name__ == "__main__":
    # https://leetcode.cn/problems/minimum-time-to-visit-all-houses/
    class Solution:
        def minTotalTime(self, forward: List[int], backward: List[int], queries: List[int]) -> int:
            n = len(forward)
            C = WeightedCycle(n, forward, backward[1:] + [backward[0]])
            res = 0
            for pre, cur in pairwise(chain([0], queries)):
                res += C.dist(pre, cur)
            return res

    forward = [1, 4, 4]
    backward = [4, 1, 2]
    queries = [1, 2, 0, 2]
    assert Solution().minTotalTime(forward, backward, queries) == 12

    weights_cw = [3, 1, 4, 2, 5]  # 顺时针边权重
    weights_ccw = [3, 1, 5, 4, 2]  # 逆时针边权重
    wc = WeightedCycle(5, weights_cw, weights_ccw)

    # 测试距离计算
    assert wc.dist(0, 2) == 4, "错误: 从0到2的最短距离应为4"
    assert wc.dist_ccw(0, 2) == 11, "错误: 从0到2的逆时针距离应为11"
    assert wc.dist_cw(0, 2) == 4, "错误: 从0到2的顺时针距离应为4"

    # 测试路径查找
    assert wc.path(0, 2) == [0, 1, 2], "错误: 从0到2的路径应为[0, 1, 2]"
    assert wc.segment(0, 2) == [(0, 2)], "错误: 从0到2的路径段应为[(0, 2)]"
    assert wc.path_ccw(0, 2) == [0, 4, 3, 2], "错误: 从0到2的逆时针路径应为[0, 4, 3, 2]"
    assert wc.segment_ccw(0, 2) == [
        (0, 0),
        (4, 2),
    ], "错误: 从0到2的逆时针路径段应为[(0, 0), (4, 2)]"
    assert wc.path_cw(0, 2) == [0, 1, 2], "错误: 从0到2的顺时针路径应为[0, 1, 2]"
    assert wc.segment_cw(0, 2) == [(0, 2)], "错误: 从0到2的顺时针路径段应为[(0, 2)]"
    assert wc.path(2, 0) == [2, 1, 0], "错误: 从2到0的路径应为[2, 1, 0]"
    assert wc.segment(2, 0) == [(2, 0)], "错误: 从2到0的路径段应为[(2, 0)]"

    # 测试路径包含性
    assert not wc.on_path_ccw(0, 2, 1), "错误: 1不在从0到2的逆时针路径上"
    assert wc.on_path_ccw(0, 2, 0), "错误: 0在从0到2的逆时针路径上"
    assert wc.on_path_ccw(0, 2, 2), "错误: 2在从0到2的逆时针路径上"
    assert wc.on_path_ccw(0, 2, 3), "错误: 3在从0到2的逆时针路径上"
    assert wc.on_path_ccw(0, 2, 4), "错误: 4在从0到2的逆时针路径上"
    assert wc.on_path_cw(0, 2, 1), "错误: 1在从0到2的顺时针路径上"

    # 测试跳跃功能
    expected = [
        (0, 0, 0),
        (1, 0, 0),
        (2, 0, 4),
        (3, 1, 4),
        (4, 2, 4),
        (5, 2, 4),
        (6, 2, 3),
        (7, 2, 3),
        (8, 3, 3),
        (9, 3, 3),
        (10, 4, 3),
        (11, 4, 2),
        (12, 4, 1),
        (13, 4, 1),
        (14, 4, 1),
        (15, 0, 0),
        (16, 0, 0),
        (17, 0, 4),
        (18, 1, 4),
        (19, 2, 4),
    ]
    for step, res_cw, res_ccw in expected:
        assert wc.jump_cw(0, step) == res_cw, f"错误: 从0出发跳跃{step}步，顺时针到达应为{res_cw}"
        assert (
            wc.jump_ccw(0, step) == res_ccw
        ), f"错误: 从0出发跳跃{step}步，逆时针到达应为{res_ccw}"
