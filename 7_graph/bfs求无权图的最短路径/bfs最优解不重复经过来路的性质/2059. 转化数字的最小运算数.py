from typing import List
from collections import deque

# 1 <= nums.length <= 1000
# 0 <= start <= 1000
# start != goal
# 使 x 越过 0 <= x <= 1000 范围的运算同样可以生效，但该运算执行后将不能执行其他运算。


# bfs 时间复杂度：O(mn)
# 其中 n 为 nums 的长度，m 为可对 x 进行操作的取值范围大小。
# 广度优先搜索至多需要将 O(m) 个数值`加入队列`，
# 对于`每个加入队列的数值可能的操作种数`为 O(n) 个


class Solution:
    def minimumOperations(self, nums: List[int], start: int, goal: int) -> int:
        queue = deque([start])
        visited = [False] * 1001
        visited[start] = True
        step = 0

        while queue:
            step += 1
            for _ in range(len(queue)):
                cur = queue.popleft()
                for num in nums:
                    for next in (cur + num, cur - num, cur ^ num):
                        if next == goal:
                            return step
                        if 0 <= next <= 1000 and not visited[next]:
                            visited[next] = True
                            queue.append(next)
        return -1

