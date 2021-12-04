from typing import List

# 假设你有一个长度为 n 的数组，初始情况下所有的数字均为 0，你将会被给出 k​​​​​​​ 个更新的操作。
# 请你返回 k 次操作后的数组。


class Solution:
    def getModifiedArray(self, length: int, updates: List[List[int]]) -> List[int]:
        res = [0] * (length + 1)
        for left, right, delta in updates:
            res[left] += delta
            res[right + 1] -= delta
        for i in range(1, length + 1):
            res[i] += res[i - 1]
        return res[:-1]


print(Solution().getModifiedArray(length=5, updates=[[1, 3, 2], [2, 4, 3], [0, 2, -2]]))
# 输出: [-2,0,3,5,3]
