from heapq import heappop, heappush, heapify
from typing import List

# 每行和每列元素均按升序排序 => 从左下角开始寻找，二分
# 每行均按升序排序，列不一定升序排序 => 每行多路归并


class Solution:
    def kthSmallest(self, matrix: List[List[int]], k: int) -> int:
        n = len(matrix)
        pq = [(matrix[i][0], i, 0) for i in range(n)]
        heapify(pq)

        for _ in range(k - 1):
            _, x, y = heappop(pq)
            if y != n - 1:
                heappush(pq, (matrix[x][y + 1], x, y + 1))

        return heappop(pq)[0]

    def kthSmallest2(self, matrix: List[List[int]], k: int) -> int:
        n = len(matrix)

        def countNGT(mid: int) -> int:
            row, col = n - 1, 0
            res = 0
            while row >= 0 and col < n:
                if matrix[row][col] <= mid:
                    res += row + 1
                    col += 1
                else:
                    row -= 1
            return res

        left, right = matrix[0][0], matrix[-1][-1]
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1

        return left


print(Solution().kthSmallest2(matrix=[[1, 5, 9], [10, 11, 13], [12, 13, 15]], k=8))
