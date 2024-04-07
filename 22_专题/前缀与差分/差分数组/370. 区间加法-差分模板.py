# 370. 区间加法
# https://leetcode.cn/problems/range-addition/description/
# 假设你有一个长度为 n 的数组，初始情况下所有的数字均为 0，你将会被给出 k​​​​​​​ 个更新的操作。
# 请你返回 k 次操作后的数组。
from typing import List
from Diff import DiffArray


class Solution:
    def getModifiedArray(self, length: int, updates: List[List[int]]) -> List[int]:
        diff = DiffArray(length)
        for left, right, delta in updates:
            diff.add(left, right + 1, delta)
        return diff.getAll()


print(Solution().getModifiedArray(length=5, updates=[[1, 3, 2], [2, 4, 3], [0, 2, -2]]))
# 输出: [-2,0,3,5,3]
