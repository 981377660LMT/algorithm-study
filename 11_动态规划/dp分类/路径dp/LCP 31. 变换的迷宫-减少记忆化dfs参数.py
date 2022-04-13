# 迷宫地形会随时间变化而改变，迷宫出口一直位于 (n-1,m-1) 位置。
# 迷宫变化规律记录于 maze 中，maze[i] 表示 i 时刻迷宫的地形状态，
# "." 表示可通行空地，"#" 表示陷阱。


from functools import lru_cache
from typing import List

# 1 <= maze.length <= 100
# 1 <= maze[i].length, maze[i][j].length <= 50
# maze[i][j] 仅包含 "."、"#"

# 记忆化dfs时间复杂度O(100*50*50*2*2)

# 临时卷轴：临时卷轴是否使用过了，由于其值仅为False或True，时间复杂度只需再乘以2。
# 永久卷轴：若直接引入使用的坐标点，则时间复杂度直接爆炸
DIRS = [(1, 0), (0, 1), (-1, 0), (0, -1), (0, 0)]


class Solution:
    def escapeMaze(self, maze: List[List[str]]) -> bool:
        @lru_cache(None)
        def dfs(index: int, x: int, y: int, isAUsed: bool, isBUsed: bool) -> bool:
            if x == row - 1 and y == col - 1:
                return True
            if index + 1 >= depth:
                return False

            for dx, dy in DIRS:
                nx, ny = x + dx, y + dy
                if 0 <= nx < row and 0 <= ny < col:
                    # 平地
                    if maze[index + 1][nx][ny] == '.':
                        if dfs(index + 1, nx, ny, isAUsed, isBUsed):
                            return True
                    # 墙
                    else:
                        if not isAUsed:
                            if dfs(index + 1, nx, ny, True, isBUsed):
                                return True
                        if not isBUsed:
                            # 对永久卷轴的等效处理技巧：保持不动，因为这个的作用就是让自己之后回来
                            for nIndex in range(index + 1, depth):
                                if dfs(nIndex, nx, ny, isAUsed, True):
                                    return True
            return False

        depth, row, col = len(maze), len(maze[0]), len(maze[0][0])
        res = dfs(0, 0, 0, False, False)
        dfs.cache_clear()
        return res


print(
    Solution().escapeMaze(
        [
            [".##..####", ".#######."],
            ["..######.", "########."],
            [".#####.##", ".#######."],
            [".#..###.#", "########."],
            [".########", "########."],
            [".######.#", "####.###."],
            [".#####.##", "#####.#.."],
            [".##.####.", "##.#####."],
            [".########", "#####.##."],
            [".#.######", "#.##.###."],
            [".########", "###.#.#.."],
            [".########", "########."],
            [".####.##.", "##.##...."],
            [".#######.", "###.#.##."],
            [".####.###", "###.####."],
            [".######.#", "##.####.."],
            [".##.#####", "##.###.#."],
            [".####.###", "##.#####."],
            [".##.##..#", ".#.#####."],
            [".###.####", "##.#..##."],
            [".####.#.#", "##.#####."],
            [".####.###", "####.###."],
            [".########", "#######.."],
            [".#####.##", "#.######."],
            [".########", "###..#.#."],
            [".####.#.#", "###..##.."],
            [".######.#", "########."],
            [".########", "##.#####."],
            [".########", "..######."],
            [".#####..#", "#######.."],
            [".#.######", ".#######."],
            [".###.#.#.", ".##..#.#."],
            [".#.##.###", "####.##.."],
        ]
    )
)

