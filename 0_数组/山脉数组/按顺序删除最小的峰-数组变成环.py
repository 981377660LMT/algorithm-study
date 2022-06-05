# 峰是比两边都大的数
# 按顺序删除最小的峰


from heapq import heapify, heappop, heappush
from typing import List


class Solution:
    def solve(self, nums: List[int]) -> List[int]:
        n = len(nums)
        # 让数组变成一个环，比较容易解决峰和谷的问题
        nums.append(-int(1e20))
        neighbors = {nums[i - 1]: [nums[i - 2], nums[i]] for i in range(n + 1)}  # 模拟的链表
        peaks = [val for val in nums if val > max(neighbors[val])]
        heapify(peaks)

        res = []
        while len(res) < n:
            peak = heappop(peaks)
            res.append(peak)
            left, right = neighbors[peak]
            neighbors[left][1] = right
            neighbors[right][0] = left
            for cand in (left, right):
                if cand > max(neighbors[cand]):
                    heappush(peaks, cand)
        return res


print(Solution().solve(nums=[3, 5, 1, 4, 2]))

# [4, 2, 5, 3, 1]
# We remove 4 and get [3, 5, 1, 2]
# We remove 2 and get [3, 5, 1]
# We remove 5 and get [3, 1]
# We remove 3 and get [1]
# We remove 1 and get []
