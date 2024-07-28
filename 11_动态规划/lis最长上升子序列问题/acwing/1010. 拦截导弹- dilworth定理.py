# dilworth定理(偏序集上最小链划分中链的数量等于其反链长度的最大值)在lis问题中的应用
# !一个序列最少用多少个不上升子序列覆盖 => 最长上升子序列的长度，他们两个是对偶问题
# !一个序列最少用多少个上升子序列覆盖 => 最长不上升子序列的长度
#
# 1010. 拦截导弹- dilworth定理
# https://www.acwing.com/problem/content/description/1012/
# 雷达给出的高度数据是不大于 30000 的正整数，导弹数不超过  1000 。
#
# 证明
# 1、首先我们把这些导弹分为s组（s即为所求答案）
# 可以看出每一组都是一个不升子序列
# 2、划分完后我们在组一里找一个原序列里以组一的开头点连续的不升子串的最后一个元素，可以知道在组2中`一定有一个大与它的点`
# （如果组二中没有的话，那么组二中最高的导弹高度必然小于这个点，而其他的高度都小于这个高度而且是递减或相等的，那么没有必要再开一个组二了，矛盾，所以不存在找不到比他大的点的情况）
# 3、以此类推，对于每一个k组（1<=k<n）都可以找到这样的一些点
# 所以把这些点连起来，就是一条上升子序列。


from typing import List, Tuple
from getLis import getLis


def 拦截导弹(nums: List[int]) -> Tuple[int, int]:
    # 第一行包含一个整数，表示最多能拦截的导弹数
    # !第二行包含一个整数，表示要拦截所有导弹最少要配备的系统数
    lis1 = getLis(nums[::-1], strict=False)[0]  # 最长不上升子序列的最长长度
    lis2 = getLis(nums)[0]  # 最少用多少个不上升子序列覆盖 => 最长上升子序列
    return len(lis1), len(lis2)


if __name__ == "__main__":

    class Solution:
        # 3231. 要删除的递增子序列的最小数量
        # https://leetcode.cn/problems/minimum-number-of-increasing-subsequence-to-be-removed/description/
        # !单调递增子序列最少划分 等于 最长不增子序列的长度
        def minOperations(self, nums: List[int]) -> int:
            lis = getLis(nums[::-1], strict=False)[0]
            return len(lis)
