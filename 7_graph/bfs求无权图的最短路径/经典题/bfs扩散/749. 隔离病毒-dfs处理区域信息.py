# 1 <= m, n <= 50
# 每天晚上，病毒会从被感染区域向相邻未感染区域扩散，除非被防火墙隔离。
# 现由于资源有限，`每天你只能安装一系列防火墙来隔离其中一个被病毒感染的区域`（一个区域或连续的一片区域），
# 且该感染区域对未感染区域的威胁最大且 保证唯一 。

# 你需要努力使得最后有部分区域不被病毒感染，如果可以成功，那么返回需要使用的防火墙个数;
# 如果无法实现，则返回在世界被病毒全部感染时已安装的防火墙个数。


from collections import defaultdict
from typing import Generator, List, Tuple

# 1 <= m, n <= 50


class Solution:
    def containVirus(self, isInfected: List[List[int]]) -> int:
        def genNexts(r: int, c: int) -> Generator[Tuple[int, int], None, None]:
            for nr, nc in [(r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)]:
                if 0 <= nr < row and 0 <= nc < col:
                    yield nr, nc

        def dfs(sr: int, sc: int) -> None:
            """dfs floodfill 处理区域信息"""
            if (sr, sc) in visited:
                return
            visited.add((sr, sc))
            region[dfsId].add((sr, sc))
            for nr, nc in genNexts(sr, sc):
                if isInfected[nr][nc] == 1:
                    dfs(nr, nc)
                elif isInfected[nr][nc] == 0:
                    border[dfsId].add((nr, nc))
                    premiter[dfsId] += 1

        row, col = len(isInfected), len(isInfected[0])
        res = 0
        dfsId = 0
        while True:
            visited = set()
            region = defaultdict(set)  # 病毒的区域
            border = defaultdict(set)  # 周围的区域
            premiter = defaultdict(int)  # 区域的周长

            # 1. 处理各个区域
            for r in range(row):
                for c in range(col):
                    if isInfected[r][c] == 1 and (r, c) not in visited:
                        dfs(r, c)
                        dfsId += 1

            if not region or not border:
                break

            # 2. 找到最大的周长
            max_ = max(map(len, border.values()))
            maxId = next((k for k, v in border.items() if len(v) == max_), -1)
            assert maxId != -1
            res += premiter[maxId]

            # 3. 所有未隔离区域扩散
            for id, points in region.items():
                if id == maxId:
                    for r, c in points:
                        isInfected[r][c] = -1
                else:
                    for r, c in points:
                        for nr, nc in genNexts(r, c):
                            if isInfected[nr][nc] == 0:
                                isInfected[nr][nc] = 1

        return res


print(
    Solution().containVirus(
        isInfected=[
            [0, 1, 0, 0, 0, 0, 0, 1],
            [0, 1, 0, 0, 0, 0, 0, 1],
            [0, 0, 0, 0, 0, 0, 0, 1],
            [0, 0, 0, 0, 0, 0, 0, 0],
        ]
    )
)

print(
    Solution().containVirus(
        isInfected=[
            [1, 1, 1, 0, 0, 0, 0, 0, 0],
            [1, 0, 1, 0, 1, 1, 1, 1, 1],
            [1, 1, 1, 0, 0, 0, 0, 0, 0],
        ]
    )
)

