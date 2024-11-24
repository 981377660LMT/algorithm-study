from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maxCollectedFruits(self, fruits: List[List[int]]) -> int:
        n = len(fruits)
        from collections import defaultdict
        from itertools import permutations

        moves = {
            1: [(1, 0), (0, 1), (1, 1)],
            2: [(1, -1), (1, 0), (1, 1)],
            3: [(-1, 1), (0, 1), (1, 1)],
        }

        starts = {
            1: (0, 0),
            2: (0, n - 1),
            3: (n - 1, 0),
        }

        def dp_kid(kid_num, grid):
            n = len(grid)
            start_i, start_j = starts[kid_num]
            dp = [{} for _ in range(n)]
            dp[0][(start_i, start_j)] = grid[start_i][start_j]

            for step in range(n - 1):
                dp_next = {}
                for (i, j), total in dp[step].items():
                    for di, dj in moves[kid_num]:
                        ni, nj = i + di, j + dj
                        if 0 <= ni < n and 0 <= nj < n:
                            new_total = total + grid[ni][nj]
                            if (ni, nj) not in dp_next or dp_next[(ni, nj)] < new_total:
                                dp_next[(ni, nj)] = new_total
                dp[step + 1] = dp_next

            max_total = 0
            if (n - 1, n - 1) in dp[n - 1]:
                max_total = dp[n - 1][(n - 1, n - 1)]
                path = set()
                dp_prev = [{} for _ in range(n)]
                dp[0][(start_i, start_j)] = grid[start_i][start_j]
                dp_prev[0][(start_i, start_j)] = None

                for step in range(n - 1):
                    dp_next = {}
                    dp_prev_next = {}
                    for (i, j), total in dp[step].items():
                        for di, dj in moves[kid_num]:
                            ni, nj = i + di, j + dj
                            if 0 <= ni < n and 0 <= nj < n:
                                new_total = total + grid[ni][nj]
                                if (ni, nj) not in dp_next or dp_next[(ni, nj)] < new_total:
                                    dp_next[(ni, nj)] = new_total
                                    dp_prev_next[(ni, nj)] = (i, j)
                    dp[step + 1] = dp_next
                    dp_prev[step + 1] = dp_prev_next

                path = set()
                i, j = n - 1, n - 1
                path.add((i, j))
                for step in range(n - 1, 0, -1):
                    i, j = dp_prev[step][(i, j)]
                    path.add((i, j))
                return max_total, path
            else:
                return 0, set()

        max_total_fruits = 0
        from itertools import permutations

        for perm in permutations([1, 2, 3]):
            grid_copy = [row[:] for row in fruits]
            total_fruits = 0
            for kid_num in perm:
                total, path = dp_kid(kid_num, grid_copy)
                total_fruits += total
                for i, j in path:
                    grid_copy[i][j] = 0
            if total_fruits > max_total_fruits:
                max_total_fruits = total_fruits

        return max_total_fruits
