from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def buildTransferStation(self, area: List[List[int]]) -> int:
        """296. 最佳的碰头地点-二维曼哈顿距离和最小"""
        ROW, COL = len(area), len(area[0])
        px = []
        py = []
        for r in range(ROW):
            for c in range(COL):
                if area[r][c] == 1:
                    px.append(c)
                    py.append(r)

        px.sort()
        py.sort()
        mid1 = px[len(px) // 2]
        mid2 = py[len(py) // 2]

        res = 0
        res += sum(abs(x - mid1) for x in px)
        res += sum(abs(y - mid2) for y in py)
        return res
