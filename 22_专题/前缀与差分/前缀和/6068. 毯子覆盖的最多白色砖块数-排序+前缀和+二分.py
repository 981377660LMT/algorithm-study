from typing import List
from bisect import bisect_right


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumWhiteTiles(self, tiles: List[List[int]], carpetLen: int) -> int:
        """枚举毯子的左端点和哪组瓷砖的左端点一致，并通过二分确定右边的位置"""
        n = len(tiles)
        tiles.sort()

        ends = [right for _, right in tiles]
        preSum = [0] * (n + 1)
        for i in range(n):
            preSum[i + 1] = preSum[i] + (tiles[i][1] - tiles[i][0] + 1)

        res = 0
        for i, (start, _) in enumerate(tiles):
            end = start + carpetLen - 1
            pos = bisect_right(ends, end)  # 右端点
            cand = preSum[pos] - preSum[i]  # 完全覆盖的
            if pos < n:
                cand += max(0, end - tiles[pos][0] + 1)  # 覆盖多的
            res = max(res, cand)

        return res


print(Solution().maximumWhiteTiles([[1, 5], [10, 11], [12, 18], [20, 25], [30, 32]], 10))

