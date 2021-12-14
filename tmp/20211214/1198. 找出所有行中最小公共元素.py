from typing import List
from collections import Counter
from functools import reduce

# 给你一个矩阵 mat，其中每一行的元素都已经按 严格递增 顺序排好了。请你帮忙找出在所有这些行中 最小的公共元素。
# 如果矩阵中没有这样的公共元素，就请返回 -1。
# 1 <= mat.length, mat[i].length <= 500


# 1. 计数排序
# 2. 集合与
# 3. 排序，可用二分:遍历第一行所有元素，然后在其余所有行使用二分搜索检查是否存在该元素。
class Solution:
    def smallestCommonElement1(self, mat: List[List[int]]) -> int:
        counter = Counter()
        for row in mat:
            counter += Counter(row)
        for num in mat[0]:
            if counter[num] == len(mat):
                return num
        return -1

    def smallestCommonElement2(self, mat: List[List[int]]) -> int:
        res = reduce(lambda i, j: i & j, [set(v) for v in mat])
        return min(res) if res else -1

    def smallestCommonElement(self, mat: List[List[int]]) -> int:
        counter = Counter()
        for row in mat:
            counter += Counter(row)
        for num in mat[0]:
            if counter[num] == len(mat):
                return num
        return -1


print(
    Solution().smallestCommonElement(
        mat=[[1, 2, 3, 4, 5], [2, 4, 5, 8, 10], [3, 5, 7, 9, 11], [1, 3, 5, 7, 9]]
    )
)
# 输出：5
