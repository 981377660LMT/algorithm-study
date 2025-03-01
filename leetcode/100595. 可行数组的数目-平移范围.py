# 100595. 可行数组的数目
# https://leetcode.cn/contest/biweekly-contest-151/problems/find-the-number-of-copy-arrays/
# 给你一个长度为 n 的数组 original 和一个长度为 n x 2 的二维数组 bounds，其中 bounds[i] = [ui, vi]。
#
# 你需要找到长度为 n 且满足以下条件的 可能的 数组 copy 的数量：
#
# 对于 1 <= i <= n - 1 ，都有 (copy[i] - copy[i - 1]) == (original[i] - original[i - 1]) 。
# 对于 0 <= i <= n - 1 ，都有 ui <= copy[i] <= vi 。
# 返回满足这些条件的数组数目。
#
# !所有满足条件的数组必须保持相同的相邻元素差值，实际上是原始数组的平移。
# !设 copy[i] = original[i] + k，其中k是一个整数常数
# 对每个元素，必须满足 bounds[i][0] <= original[i] + k <= bounds[i][1]
# 转化为 bounds[i][0] - original[i] <= k <= bounds[i][1] - original[i]

from typing import List


INF = int(1e20)


class Solution:
    def countArrays(self, original: List[int], bounds: List[List[int]]) -> int:
        min_, max_ = -INF, INF
        for i in range(len(original)):
            min_ = max(min_, bounds[i][0] - original[i])
            max_ = min(max_, bounds[i][1] - original[i])
        return max(0, max_ - min_ + 1)
