from collections import deque

# 骑士需要前去征服坐标为 [x, y] 的部落，请你为他规划路线。

# 最后返回所需的最小移动次数即可。本题确保答案是一定存在的。

# 1.4个象限关于（0，0）对称。都转化成第一象限，方便处理和计算
# 2.限定nextX,nextY的范围剪枝


class Solution:
    def minKnightMoves(self, x: int, y: int) -> int:
        if x == 0 and y == 0:
            return 0
        x = abs(x)  # 关于（0,0）对称
        y = abs(y)

        # --------------------------------记忆化 bfs------------------------------
        queue = deque([(0, 0)])
        visited = set([(0, 0)])

        step = 0
        while queue:
            cur_len = len(queue)
            step += 1
            for _ in range(cur_len):
                [x0, y0] = queue.popleft()
                for dx, dy in (
                    (-2, 1),
                    (-2, -1),
                    (-1, 2),
                    (-1, -2),
                    (1, 2),
                    (1, -2),
                    (2, 1),
                    (2, -1),
                ):
                    nx = x0 + dx
                    ny = y0 + dy
                    if (nx, ny) not in visited:
                        # 路径不能是往左10再往右10，往左100，再往右100. 相当于自己做了一个评估函数剪枝剪枝
                        if -5 <= nx <= x + 5 and -5 <= ny <= y + 5:
                            if nx == x and ny == y:
                                return step
                            else:
                                queue.append((nx, ny))
                                visited.add((nx, ny))


print(Solution().minKnightMoves(x=5, y=5))
