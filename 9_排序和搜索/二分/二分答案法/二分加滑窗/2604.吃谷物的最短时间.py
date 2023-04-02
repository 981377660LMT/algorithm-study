# 吃谷物的最短时间
# 在一条线上有 n 个母鸡和 m 个谷粒。
# 给出了两个整数阵列中母鸡和谷物的初始位置，母鸡和谷物的大小分别为 n 和 m。
# 任何母鸡都可以吃谷物，如果他们在同一位置。花在这上面的时间可以忽略不计。
# 一只母鸡也可以吃多种谷物。在1秒内，母鸡可以向左或向右移动1个单位。
# `母鸡可以同时独立地移动`。如果母鸡表现最佳，返回吃所有谷物的最短时间。
# n<=2e4
# !0<=pos<=1e9

# !注意不是最短移动距离, 而是总时间(并行)

# 二分+排序
# 1.每只母鸡应该按顺序匹配谷物，因此我们可以对母鸡和谷物的位置进行排序。
# 2.二分答案,每只鸡在给定的时间内可以吃多少粒谷物
# !3.调头一次的最短距离=min(2*leftMax+rightMax,2*rightMax+leftMax)
# !其中leftMax为向左走的最大距离,rightMax为向右走的最大距离
# (摘水果那道题也是这个调头公式)

from typing import List


class Solution:
    def minimumTime(self, hens: List[int], grains: List[int]) -> int:
        def calDist(start: int, left: int, right: int) -> int:
            """从start出发,遍历[leftMost,rightMost]区间的最短距离(最多调头一次)"""
            leftMax = max(0, start - left)
            rightMax = max(0, right - start)
            return min(2 * leftMax + rightMax, 2 * rightMax + leftMax)

        def check(mid: int) -> bool:
            left, eat = 0, 0
            for hStart in hens:
                right = left
                while right < len(grains) and calDist(hStart, grains[left], grains[right]) <= mid:
                    right += 1
                eat += right - left
                left = right
            return eat == len(grains)

        hens.sort()
        grains.sort()
        left, right = 0, int(1e12)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


assert Solution().minimumTime([3], [5, 1, 3, 9, 3]) == 10
