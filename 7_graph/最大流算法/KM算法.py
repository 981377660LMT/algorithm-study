class KM:
    INF = int(1e9)

    def __init__(self, graph):
        ma = max(len(graph), len(graph[0]))
        self.graph = graph
        self.vis_x = set()
        self.vis_y = set()
        self.match = [-1] * ma
        self.lx = [max(row) for row in graph]
        self.ly = [0] * ma
        self.slack = []

        self.nx, self.ny = len(graph), len(graph[0])

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
        ans = 0
        # print(KM(g).solve())
        for i, j in enumerate(KM(g).solve()):
            if j != -1:
                ans += g[j][i]
        return ans
