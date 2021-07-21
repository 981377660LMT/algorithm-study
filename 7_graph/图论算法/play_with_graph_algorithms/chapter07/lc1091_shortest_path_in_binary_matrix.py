from collections import deque


class Solution:

    dirs = [
        [-1, 0],
        [-1, 1],
        [0, 1],
        [1, 1],
        [1, 0],
        [1, -1],
        [0, -1],
        [-1, -1],
    ]

    def shortest_path_binary_matrix(self, grid):
        if not grid or not grid[0]:
            return -1
            
        self._R, self._C = len(grid), len(grid[0])
        visited = [[False] * self._C for _ in range(self._R)]
        dis = [[0] * self._C for _ in range(self._R)]

        if grid[0][0] == 1:
            return -1

        if self._R == 0 and self._C == 0:
            return 1

        queue = deque()
        queue.append(0)
        visited[0][0] = True
        # dis saves how many optimized points to get here
        dis[0][0] = 1

        while queue:
            cur = queue.popleft()
            curx, cury = cur // self._C, cur % self._C
            for dx, dy in self.dirs:
                nextx = curx + dx
                nexty = cury + dy
                if self._in_area(nextx, nexty) and not visited[nextx][nexty] and grid[nextx][nexty] == 0:
                    queue.append(nextx * self._C + nexty)
                    visited[nextx][nexty] = True
                    dis[nextx][nexty] = dis[curx][cury] + 1
                    if nextx == self._R - 1 and nexty == self._C - 1:
                        return dis[nextx][nexty]
        
        return -1

    
    def _in_area(self, x, y):
        return 0 <= x < self._R and 0 <= y < self._C


if __name__ == '__main__':
    data = [[0,0,0],[1,1,0],[1,1,0]]
    sol = Solution()
    print(sol.shortest_path_binary_matrix(data))
