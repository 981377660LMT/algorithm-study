from typing import List


# 3 <= A.length <= 10000
# 两边之和大于第三边是充要条件

# 排序；看最大的边
class Solution:
    def largestPerimeter(self, nums: List[int]) -> int:
        edges = sorted(nums, reverse=True)
        return next((x + y + z for x, y, z in zip(edges, edges[1:], edges[2:]) if x < y + z), 0)

    def largestPerimeter2(self, nums: List[int]) -> int:
        edges = sorted(nums, reverse=True)
        for i in range(len(edges) - 2):
            if edges[i] < edges[i + 1] + edges[i + 2]:
                return edges[i] + edges[i + 1] + edges[i + 2]
        return 0


print(Solution().largestPerimeter([3, 2, 3, 4]))
