# 和为n的最少完全平方数
# 1 ≤ n ≤ 100,000

# 时间复杂度O(nsqrtn)
from collections import deque


class Solution:
    def solve(self, n):
        squares = []
        for i in range(1, int(n ** 0.5) + 1):
            squares.append(i ** 2)
        squares = squares[::-1]   # 倒着搜索更快

        queue = deque([0])
        visited = set()
        res = 0

        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                if cur == n:
                    return res
                if cur > n:
                    continue

                if cur in visited:
                    continue
                visited.add(cur)

                for add in squares:
                    queue.append(cur + add)

            res += 1

        return res
