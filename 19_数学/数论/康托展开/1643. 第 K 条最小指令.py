from typing import List
from 康托展开 import calPerm

# 他只能向 右 或向 下 走。你可以为 Bob 提供导航 指令 来帮助他到达目的地 destination 。
# 指令 用字符串表示，其中每个字符：
# 'H' ，意味着水平向右移动
# 'V' ，意味着竖直向下移动
# 给你一个整数数组 destination 和一个整数 k ，请你返回可以为 Bob 提供前往目的地 destination 导航的
# 按字典序排列后的第 k 条最小指令 。
# 1 <= row, column <= 15
# 总结:求字典序第k小的排列


class Solution:
    def kthSmallestPath(self, destination: List[int], k: int) -> str:
        res = calPerm(["H"] * destination[1] + ["V"] * destination[0], k - 1)
        return "".join(res)


print(Solution().kthSmallestPath(destination=[2, 3], k=1))
# 输出："HHHVV"
# 解释：能前往 (2, 3) 的所有导航指令 按字典序排列后 如下所示：
# ["HHHVV", "HHVHV", "HHVVH", "HVHHV", "HVHVH", "HVVHH", "VHHHV", "VHHVH", "VHVHH", "VVHHH"].
