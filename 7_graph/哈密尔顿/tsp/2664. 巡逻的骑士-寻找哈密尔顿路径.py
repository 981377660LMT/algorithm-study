# 给定两个正整数 m 和 n ，它们是一个 下标从 0 开始 的二维数组 board 的高度和宽度。
# 还有一对正整数 (r, c) ，它们是骑士在棋盘上的起始位置。
# 你的任务是找到一个骑士的移动顺序，使得 board 中每个单元格都 恰好 被访问一次（起始单元格已被访问，不应 再次访问）。
# 返回数组 board ，其中单元格的值显示从 0 开始访问该单元格的顺序（骑士的初始位置为 0）。


# 1 <= m, n <= 5
# 0 <= r <= m - 1
# 0 <= c <= n - 1
# 输入的数据保证在给定条件下至少存在一种访问所有单元格的移动顺序。

# 2664. 巡逻的骑士-回溯法寻找哈密尔顿路径


from typing import List


DIR8 = ((-2, -1), (-2, 1), (-1, -2), (-1, 2), (1, -2), (1, 2), (2, -1), (2, 1))


class Solution:
    def tourOfKnight(self, m: int, n: int, r: int, c: int) -> List[List[int]]:
        def bt(x: int, y: int):
            if len(path) == ok:
                yield path
                return
            for dx, dy in DIR8:
                nx, ny = x + dx, y + dy
                if 0 <= nx < m and 0 <= ny < n and not visited[nx][ny]:
                    visited[nx][ny] = True
                    path.append((nx, ny))
                    yield from bt(nx, ny)
                    path.pop()
                    visited[nx][ny] = False

        ok = m * n
        path = [(r, c)]
        visited = [[False] * n for _ in range(m)]
        visited[r][c] = True
        resPath = next(bt(r, c))
        res = [[0] * n for _ in range(m)]
        for i, (x, y) in enumerate(resPath):
            res[x][y] = i
        return res
