from typing import List

from sortedcontainers import SortedList


class Solution:
    def minDayskVariants(self, P: List[List[int]], k: int) -> int:
        """
        给定二维数组 points ，第 i 项 points[i] = [xi, yi] 说明第 0 天有一种病毒在点 (xi, yi) 。
        注意初始状态下，可能有 多种 病毒在 同一点 上。
        每天，被感染的点会把它感染的病毒传播到上、下、左、右四个邻居点。
        现给定一个整数 k ，问 最少 需要多少天，方能找到一点感染 至少 k 种病毒？

        2 <= n <= 50
        1 <= xi, yi <= 1e9

        https://leetcode.cn/problems/minimum-time-for-k-virus-variants-to-spread/comments/1178385
        二分+滑动窗口
        可以二分转为判断问题：在变换的坐标系下，给定时间 t，能否找到一个对角线长 2*t 的菱形，包含了 k 个点。
        注意按题意找到的正方形中心必须在网格上
        """
        # 将`坐标系`逆时针旋转 45 度(即点顺时针旋转45度)并扩大根号2倍

        def check(mid: int) -> bool:
            """能否找到一个长 2*mid 的滑窗，包含了 k 个点"""
            left = 0
            sl = SortedList()
            for x, y in points:
                sl.add(y)
                while sl and x - points[left][0] > 2 * mid:
                    sl.discard(points[left][1])
                    left += 1

                # 检查纵坐标
                for start in range(len(sl) - k + 1):
                    # 按题意找到的正方形中心必须在网格上，因此还需要加个判断 TODO
                    if sl[start + k - 1] - sl[start] <= 2 * mid:
                        return True
            return False

        points = sorted((y + x, y - x) for x, y in P)
        left, right = 0, int(1e14)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
