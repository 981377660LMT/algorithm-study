from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 n ，你可以执行下述操作 任意 次：

# n 加上或减去 2 的某个 幂
# 返回使 n 等于 0 需要执行的 最少 操作数。

# 如果 x == 2i 且其中 i >= 0 ，则数字 x 是 2 的幂。
class Solution:
    def minOperations(self, n: int) -> int:
        # 如果是2的幂 为0
        if n & (n - 1) == 0:
            return 0
        if n == 1:
            return 1
        queue = deque([(n, 0)])
        visited = set([n])
        while queue:
            cur, step = queue.popleft()
            if cur == 0:
                return step
            for i in range(20):
                nxt = cur - (1 << i)
                if nxt >= 0 and nxt not in visited:
                    queue.append((nxt, step + 1))
                    visited.add(nxt)
                nxt = cur + (1 << i)
                if nxt not in visited and nxt <= 1 << 18:
                    queue.append((nxt, step + 1))
                    visited.add(nxt)


print(1 << 18)
print(Solution().minOperations(54))
print(Solution().minOperations(39))
