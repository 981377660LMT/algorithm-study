# 100388. 交替组 III-环形前缀和查询
# https://leetcode.cn/problems/alternating-groups-iii/description/
# 给你一个整数数组 colors 和一个二维整数数组 queries 。colors表示一个由红色和蓝色瓷砖组成的环，第 i 块瓷砖的颜色为 colors[i] ：

# colors[i] == 0 表示第 i 块瓷砖的颜色是 红色 。
# colors[i] == 1 表示第 i 块瓷砖的颜色是 蓝色 。
# 环中连续若干块瓷砖的颜色如果是 交替 颜色（也就是说这组瓷砖中除了第一块和最后一块瓷砖以外，中间瓷砖的颜色与它 左边 和 右边 的颜色都不同），那么它被称为一个 交替组。

# 你需要处理两种类型的查询：

# queries[i] = [1, sizei]，确定大小为sizei的 交替组 的数量。
# queries[i] = [2, indexi, colori]，将colors[indexi]更改为colori。
# 返回数组 answer，数组中按顺序包含第一种类型查询的结果。

# 注意 ，由于 colors 表示一个 环 ，第一块 瓷砖和 最后一块 瓷砖是相邻的。

# !环形前缀和3


from typing import List


class Solution:
    def numberOfAlternatingGroups(self, colors: List[int], queries: List[List[int]]) -> List[int]:
        ...
