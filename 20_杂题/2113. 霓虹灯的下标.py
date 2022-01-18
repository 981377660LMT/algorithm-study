from typing import List

# 1 <= nums.length <= 100
# 1 <= n <= 105
# 2n 一循环；使用切片防止 IndexError
class Solution:
    def elementInNums(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        n = len(nums)
        states = []
        for i in range(n):
            states.append(nums[i:])
        for i in range(n):
            states.append(nums[:i])

        res = []
        for cycle, index in queries:
            cycle %= len(states)
            cand = states[cycle][index : index + 1]
            res.append(cand[0] if cand else -1)

        return res


print(Solution().elementInNums(nums=[0, 1, 2], queries=[[0, 2], [2, 0], [3, 2], [5, 0]]))
print(Solution().elementInNums(nums=[2], queries=[[0, 0], [1, 0], [2, 0], [3, 0]]))

