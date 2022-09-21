from typing import List
from BIT import BIT2

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumWhiteTiles(self, tiles: List[List[int]], carpetLen: int) -> int:
        """树状数组区间更新查询"""
        # bit = BIT2(int(1e9 + 10))
        bit = BIT2(max(max(row) for row in tiles) + 10)
        res = 0
        for left, right in tiles:
            bit.add(left, right, 1)
        # 毯子的左端点一定和某组瓷砖的左端点一致。
        for left, _ in tiles:
            res = max(res, bit.query(left, left + carpetLen - 1))
        return res
