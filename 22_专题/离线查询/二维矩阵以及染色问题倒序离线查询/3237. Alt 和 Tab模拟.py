# 3237. Alt 和 Tab模拟
# https://leetcode.cn/problems/alt-and-tab-simulation/description/
# 有 n 个编号从  1 到 n 的打开的窗口，我们想要模拟使用 alt + tab 键在窗口之间导航。
# 给定数组 windows 包含窗口的初始顺序（第一个元素在最前面，最后一个元素在最后面）。
# 同时给定数组 queries 表示每一次查询中，窗口 queries[i] 被切换到最前面。
# 返回 windows 数组的最后状态。
#
# 倒序遍历操作.

from typing import List


class Solution:
    def simulationResult(self, windows: List[int], queries: List[int]) -> List[int]:
        """离线."""
        visited = set()
        res = []
        for v in queries[::-1]:
            if v not in visited:
                visited.add(v)
                res.append(v)
        for v in windows:
            if v not in visited:
                res.append(v)
        return res

    def simulationResult2(self, windows: List[int], queries: List[int]) -> List[int]:
        """在线."""
        res = {v: True for v in windows[::-1]}  # 插入有序
        for v in queries:
            del res[v]
            res[v] = True
        return list(res)[::-1]
