# 1503. 所有蚂蚁掉下来前的最后一刻
# https://leetcode.cn/problems/last-moment-before-all-ants-fall-out-of-a-plank/
# 一些蚂蚁在木板上移动，每只蚂蚁都以 每秒一个单位 的速度移动。其中，一部分蚂蚁向 左 移动，其他蚂蚁向 右 移动。
# 当两只向 不同 方向移动的蚂蚁在某个点相遇时，它们会同时改变移动方向并继续移动。假设更改方向不会花费任何额外时间。
# 两个数组分别标识向左或者向右移动的蚂蚁在 t = 0 时的位置。
# !请你返回最后一只蚂蚁从木板上掉下来的时刻。
#
# 类似于：一条狗在两个人之间来回走，问狗最后走了多少米
#
# !两只相遇的蚂蚁同时改变移动方向之后的情形等价于两只蚂蚁都不改变移动方向
# (因为我们区分不出哪知蚂蚁是哪只，相当于没有交换)


from typing import List


class Solution:
    def getLastMoment(self, n: int, left: List[int], right: List[int]) -> int:
        leftMax = max(left, default=0)
        rightMax = n - min(right, default=n)
        return max(leftMax, rightMax)
