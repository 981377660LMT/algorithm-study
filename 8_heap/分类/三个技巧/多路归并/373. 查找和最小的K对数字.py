from typing import List
from heapq import heappush, heappop


class Solution:
    def kSmallestPairs(self, nums1: List[int], nums2: List[int], k: int) -> List[List[int]]:
        """给定两个以升序排列的整形数组 nums1 和 nums2, 以及一个整数 k。
        定义一对值 (u,v),其中第一个元素来自 nums1,第二个元素来自 nums2。
        找到和最小的 k 对数字 (u1,v1), (u2,v2) ... (uk,vk)。
        """

        n1, n2 = len(nums1), len(nums2)
        res = []
        pq = [(nums1[i] + nums2[0], i, 0) for i in range(n1)]
        while len(res) < k and pq:
            _, row, col = heappop(pq)
            res.append([nums1[row], nums2[col]])
            if col + 1 < n2:
                heappush(pq, (nums1[row] + nums2[col + 1], row, col + 1))
        return res


print(Solution().kSmallestPairs(nums1=[1, 7, 11], nums2=[2, 4, 6], k=3))

# 输出: [1,2],[1,4],[1,6]
# 解释: 返回序列中的前 3 对数：
#      [1,2],[1,4],[1,6],[7,2],[7,4],[11,2],[7,6],[11,4],[11,6]
