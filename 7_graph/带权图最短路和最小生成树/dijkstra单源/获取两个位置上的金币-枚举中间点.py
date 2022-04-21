import heapq

# 2203. 得到要求路径的最小带权子图-正反图+枚举中间点

# 从起点出发收集到两个金币的最短路径
# 枚举中间点即可
class Solution:
    def solve(self, matrix, row, col, erow0, ecol0, erow1, ecol1):
        def isValid(x, y):
            return x >= 0 and y >= 0 and x < len(matrix) and y < len(matrix[0])

        def dijk(sx, sy):
            heap = [(matrix[sx][sy], sx, sy)]
            dists = [[int(1e20)] * len(matrix[0]) for _ in range(len(matrix))]
            dists[sx][sy] = matrix[sx][sy]
            while heap:
                cost, x, y = heapq.heappop(heap)
                if dists[x][y] < cost:
                    continue
                for nx, ny in [(x, y - 1), (x + 1, y), (x - 1, y), (x, y + 1)]:
                    if isValid(nx, ny) and matrix[nx][ny] + cost < dists[nx][ny]:
                        weight = matrix[nx][ny]
                        dists[nx][ny] = weight + cost
                        heapq.heappush(heap, (weight + cost, nx, ny))
            return dists

        res = int(1e20)
        # 三个点的最短路
        a, b, c = dijk(row, col), dijk(erow0, ecol0), dijk(erow1, ecol1)
        for i in range(len(matrix)):
            for j in range(len(matrix[0])):
                # 枚举中间点
                res = min(res, a[i][j] + b[i][j] + c[i][j] - 2 * matrix[i][j])
        return res
