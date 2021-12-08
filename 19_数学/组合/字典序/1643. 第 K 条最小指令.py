from typing import List
from math import comb

# 他只能向 右 或向 下 走。你可以为 Bob 提供导航 指令 来帮助他到达目的地 destination 。
# 指令 用字符串表示，其中每个字符：

# 'H' ，意味着水平向右移动
# 'V' ，意味着竖直向下移动

# 给你一个整数数组 destination 和一个整数 k ，请你返回可以为 Bob 提供前往目的地 destination 导航的
# 按字典序排列后的第 k 条最小指令 。

# 1 <= row, column <= 15

# 总结:求第k个组合
class Solution:
    def kthSmallestPath(self, destination: List[int], k: int) -> str:
        # 本题为组合题，有h个相同的'H'和v个相同的'V',我们需要对他们进行排序
        # 队列长度为h+v,我们需要选择h个位置给'H'
        # 也就是说对h,v来说，组合数有comb(h+v,h)个
        v, h = destination
        res = ''
        while h and v:
            #  此处放'H'的组合数
            num = comb(h + v - 1, h - 1)
            # 如果放'H'都不足以凑够k个，说明这个字典序太小了，这里要放V，然后在后面找k-num
            if k > num:
                res += 'V'
                v -= 1
                k -= num
            else:
                res += 'H'
                h -= 1
        return res + h * 'H' + v * 'V'


print(Solution().kthSmallestPath(destination=[2, 3], k=1))
# 输出："HHHVV"
# 解释：能前往 (2, 3) 的所有导航指令 按字典序排列后 如下所示：
# ["HHHVV", "HHVHV", "HHVVH", "HVHHV", "HVHVH", "HVVHH", "VHHHV", "VHHVH", "VHVHH", "VVHHH"].

