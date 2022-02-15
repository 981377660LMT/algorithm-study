# https://www.bilibili.com/video/BV16K4y1X7Ph/?spm_id_from=333.788.recommend_more_video.0

# O(n^3)
from typing import List


class KM:
    """注意建图时让二分图两半相等"""

    INF = int(1e9)

    def __init__(self, adjMatrix: List[List[int]]):
        ma = max(len(adjMatrix), len(adjMatrix[0]))
        self.graph = adjMatrix
        self.vis_x = set()
        self.vis_y = set()
        self.match = [-1] * ma
        self.lx = [max(row) for row in adjMatrix]
        self.ly = [0] * ma
        self.slack = []

        self.nx, self.ny = len(adjMatrix), len(adjMatrix[0])

    def find_path(self, x):
        self.vis_x.add(x)
        for y in range(self.ny):
            if y in self.vis_y:
                continue
            tmp_delta = self.lx[x] + self.ly[y] - self.graph[x][y]
            if tmp_delta == 0:
                self.vis_y.add(y)
                if self.match[y] == -1 or self.find_path(self.match[y]):
                    self.match[y] = x
                    return True
            elif self.slack[y] > tmp_delta:
                self.slack[y] = tmp_delta
        return False

    def solve(self):
        for x in range(self.nx):
            self.slack = [KM.INF] * self.ny
            while True:
                self.vis_x.clear()
                self.vis_y.clear()
                if self.find_path(x):
                    break
                else:
                    delta = KM.INF
                    for j in range(self.ny):
                        if j not in self.vis_y and delta > self.slack[j]:
                            delta = self.slack[j]
                    for i in range(self.nx):
                        if i in self.vis_x:
                            self.lx[i] -= delta
                    for j in range(self.ny):
                        if j in self.vis_y:
                            self.ly[j] += delta
                        else:
                            self.slack[j] -= delta

        return self.match


class Solution:
    def maximumANDSum(self, nums: List[int], numSlots: int) -> int:
        slots = list(range(1, numSlots + 1)) + list(range(1, numSlots + 1))
        g = [[0 for _ in slots] for _ in nums]
        for i in range(len(nums)):
            for j in range(numSlots * 2):
                g[i][j] = nums[i] & slots[j]
        res = 0
        # print(KM(g).solve())
        for i, j in enumerate(KM(g).solve()):
            if j != -1:
                res += g[j][i]
        return res
