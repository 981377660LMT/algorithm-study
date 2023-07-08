# 一个长度为 n 下标从 0 开始的整数数组 arr 的 不平衡数字 定义为，
# 在 sarr = sorted(arr) 数组中，满足以下条件的下标数目：

# 0 <= i < n - 1 ，和
# sarr[i+1] - sarr[i] > 1
# 这里，sorted(arr) 表示将数组 arr 排序后得到的数组。
# 给你一个下标从 0 开始的整数数组 nums ，请你返回它所有 子数组 的 不平衡数字 之和。
# 子数组指的是一个数组中连续一段 非空 的元素序列。

# 1 <= nums.length <= 1000
# 1 <= nums[i] <= nums.length

# !枚举子数组的同时，维护不平衡度


from typing import List


class Solution:
    def sumImbalanceNumbers(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        for i, num1 in enumerate(nums):
            visited = [False] * (n + 5)
            visited[num1] = True
            cur = 0
            for j in range(i + 1, n):
                num2 = nums[j]
                if not visited[num2]:
                    visited[num2] = True
                    cur += 1 - visited[num2 - 1] - visited[num2 + 1]
                res += cur
        return res
