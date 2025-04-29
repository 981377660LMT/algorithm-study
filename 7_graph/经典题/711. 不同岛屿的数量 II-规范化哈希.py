from typing import List, Set, Tuple


class Solution:
    def numDistinctIslands2(self, grid: List[List[int]]) -> int:
        """
        统计不同岛屿的形状，允许旋转（90°、180°、270°）和镜像（水平、垂直）。
        对每个岛屿：
          1. 用 DFS 收集它的所有坐标 (i, j)。
          2. 将坐标集视作一个形状，生成其 8 种对称变换后的坐标集：
             ( x, y), ( x,-y),(-x, y),(-x,-y), ( y, x),( y,-x),(-y, x),(-y,-x)
          3. 对每个变换，平移到原点：(x - min_x, y - min_y)，排序后转为元组。
          4. 取这 8 个元组中字典序最小的一个作为该岛屿的“规范形状”。
        最后把所有岛屿的规范形状放入集合，集合大小即为不同岛屿数。
        时间 O(m·n + K·log K)，空间 O(m·n)，其中 K 是单个岛屿的格子数累加。
        """
        if not grid or not grid[0]:
            return 0
        m, n = len(grid), len(grid[0])
        seen = [[False] * n for _ in range(m)]

        def dfs(i: int, j: int, coords: List[Tuple[int, int]]):
            """从 (i,j) 向四周扩展，收集一个岛屿的所有坐标"""
            stack = [(i, j)]
            seen[i][j] = True
            while stack:
                x, y = stack.pop()
                coords.append((x, y))
                for dx, dy in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                    nx, ny = x + dx, y + dy
                    if 0 <= nx < m and 0 <= ny < n and not seen[nx][ny] and grid[nx][ny] == 1:
                        seen[nx][ny] = True
                        stack.append((nx, ny))

        def normalize(coords: List[Tuple[int, int]]) -> Tuple[Tuple[int, int], ...]:
            """对一个坐标列表，生成其规范形状元组"""
            shapes = []
            pts = coords
            # 转换为相对于原点的坐标
            # 先把所有点转换到以 (0,0) 为中心的坐标系
            for transform in (
                lambda x, y: (x, y),
                lambda x, y: (x, -y),
                lambda x, y: (-x, y),
                lambda x, y: (-x, -y),
                lambda x, y: (y, x),
                lambda x, y: (y, -x),
                lambda x, y: (-y, x),
                lambda x, y: (-y, -x),
            ):
                trans = [transform(x - base_x, y - base_y) for x, y in pts]
                # 平移到非负
                min_x = min(x for x, _ in trans)
                min_y = min(y for _, y in trans)
                shifted = sorted((x - min_x, y - min_y) for x, y in trans)
                shapes.append(tuple(shifted))
            # 取 8 种变换中字典序最小的
            return min(shapes)

        distinct: Set[Tuple[Tuple[int, int], ...]] = set()
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1 and not seen[i][j]:
                    coords: List[Tuple[int, int]] = []
                    dfs(i, j, coords)
                    base_x, base_y = coords[0]
                    distinct.add(normalize(coords))
        return len(distinct)
