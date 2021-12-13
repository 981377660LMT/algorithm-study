from typing import List

# 现在你需要从n个不同的数组中选择两个整数并且计算它们的距离|x-y|
# 你的任务就是去找到最大距离

# 总结：遍历，不断维护min与max
class Solution:
    def maxDistance(self, arrays: List[List[int]]) -> int:
        n = len(arrays)
        res = 0
        minVal = arrays[0][0]
        maxVal = arrays[0][-1]
        for i in range(1, n):
            res = max(res, abs(arrays[i][-1] - minVal), abs(arrays[i][0] - maxVal))
            minVal = min(minVal, arrays[i][0])
            maxVal = max(maxVal, arrays[i][-1])
        return res


print(Solution().maxDistance([[1, 2, 3], [4, 5], [1, 2, 3]]))
# 输出： 4
# 解释：
# 一种得到答案 4 的方法是从第一个数组或者第三个数组中选择 1，同时从第二个数组中选择 5 。
