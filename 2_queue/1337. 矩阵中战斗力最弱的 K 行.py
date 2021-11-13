# 矩阵由若干军人和平民组成，分别用 1 和 0 表示。
# 关键：军人 总是 排在一行中的靠前位置
# 也就是说 1 总是出现在 0 之前。(暗示二分统计军人数)
# 请你返回矩阵中战斗力最弱的 k 行的索引，按从最弱到最强排序(暗示固定容量k的大根堆)
from typing import List
from bisect import bisect_right
from heapq import heappush, heappushpop


class Solution:
    def kWeakestRows(self, mat: List[List[int]], k: int) -> List[int]:
        m, n = len(mat), len(mat[0])
        res = []
        for row in range(m):
            cnt = -(n - bisect_right(mat[row][::-1], 0))
            if len(res) < k:
                heappush(res, (cnt, -row))
            else:
                heappushpop(res, (cnt, -row))
        res.sort(key=lambda item: (-item[0], -item[1]))
        return [-item[1] for item in res]

