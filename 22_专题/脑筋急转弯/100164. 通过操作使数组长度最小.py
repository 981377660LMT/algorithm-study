# 100164. 通过操作使数组长度最小
# https://leetcode.cn/contest/biweekly-contest-122/problems/minimize-length-of-array-using-operations/
# 给你一个下标从 0 开始的整数数组 nums ，它只包含 正 整数。
# 你的任务是通过进行以下操作 任意次 （可以是 0 次） 最小化 nums 的长度：
# !在 nums 中选择 两个不同 的下标 i 和 j ，满足 nums[i] > 0 且 nums[j] > 0 。
# !将结果 nums[i] % nums[j] 插入 nums 的结尾。
# !将 nums 中下标为 i 和 j 的元素删除。
# !请你返回一个整数，它表示进行任意次操作以后 nums 的 最小长度 。
#
# nums[i]<=1e9, 1<=nums.length<=1e5
#
# 从样例里发现小的数可以直接吃掉大的数
# 所以只用关注最小的数有几个，以及能不能构造出一个更小的数。如果能构造出最小的数，可以只剩下这个数。
# 如果有数不是最小的数的倍数，那肯定能构造出一个更小的数；如果所有数都是它的倍数，那一定构造不出来（只能构造出0和最小的数的倍数）
# 假设最小的数的数量是 x，最后就最少剩下 (x+1)//2 个数


from collections import Counter
from typing import List


class Solution:
    def minimumArrayLength(self, nums: List[int]) -> int:
        counter = Counter(nums)
        min_ = min(nums)
        if all([x % min_ == 0 for x in counter]):
            return (counter[min_] + 1) // 2
        return 1
