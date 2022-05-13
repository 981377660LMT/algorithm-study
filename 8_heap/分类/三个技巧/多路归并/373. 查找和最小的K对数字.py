from typing import List
from heapq import heappush, heappop


class Solution:
    def kSmallestPairs(self, nums1: List[int], nums2: List[int], k: int) -> List[List[int]]:
        queue = []

        def push(i: int, j: int):
            if i < len(nums1) and j < len(nums2):
                heappush(queue, (nums1[i] + nums2[j], i, j))

        push(0, 0)
        visited = set()
        res = []

        while queue and len(res) < k:
            _, i, j = heappop(queue)
            res.append([nums1[i], nums2[j]])
            if (i, j + 1) not in visited:
                push(i, j + 1)
                visited.add((i, j + 1))
            if (i + 1, j) not in visited:
                push(i + 1, j)
                visited.add((i + 1, j))

        return res


print(Solution().kSmallestPairs(nums1=[1, 7, 11], nums2=[2, 4, 6], k=3))

# 输出: [1,2],[1,4],[1,6]
# 解释: 返回序列中的前 3 对数：
#      [1,2],[1,4],[1,6],[7,2],[7,4],[11,2],[7,6],[11,4],[11,6]

