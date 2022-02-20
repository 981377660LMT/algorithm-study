from typing import Deque, List, Tuple
from collections import deque
from itertools import combinations

# 在一个学期中，你 最多 可以同时上 k 门课，前提是这些课的先修课在之前的学期里已经上过了。
# 请你返回上完所有课最少需要多少个学期。题目保证一定存在一种上完所有课的方式。

# 1 <= n <= 15  暗示2^n 状压

# 这道题其实是说，n个用时相同的任务（work）有先后依赖关系，现只有k台机器（parallelism 并行度），最少用时多少能完成。最大深度为depth。

# bfs `最短路问题`
# https://leetcode-cn.com/problems/parallel-courses-ii/solution/bian-xing-de-bfszui-duan-lu-wen-ti-by-ca-bio2/初始化pre依赖项和dist数组


# 初始化pre依赖项和dist数组
# 初始化queue，bfs模板
# 枚举n(状压题特点) 看今年学什么课(排除调不能学的和学过的)
# 可以学的小于k 根据dist约束入队
# 可以学的大于k 选k个学


class Solution:
    def minNumberOfSemesters(self, n: int, dependencies: List[List[int]], k: int) -> int:
        pre = [0] * n
        dist = [0x7FFFFFFF] * (1 << n)
        target = (1 << n) - 1
        for u, v in dependencies:
            u, v = u - 1, v - 1
            pre[v] |= 1 << u

        queue: Deque[Tuple[int, int]] = deque([(0, 0)])  # (state,cost)
        while queue:
            state, cost = queue.popleft()
            if state == target:
                return cost

            canStudy = []
            # 学不了/学过了
            for i in range(n):
                if pre[i] & state != pre[i]:
                    continue
                if (1 << i) & state:
                    continue
                canStudy.append(i)

            if len(canStudy) <= k:
                for course in canStudy:
                    state |= 1 << course
                if dist[state] > cost + 1:
                    dist[state] = cost + 1
                    queue.append((state, cost + 1))
            else:
                for studyComb in combinations(canStudy, k):
                    nextState = state
                    for course in studyComb:
                        nextState |= 1 << course
                    if dist[nextState] > cost + 1:
                        dist[nextState] = cost + 1
                        queue.append((nextState, cost + 1))


print(Solution().minNumberOfSemesters(n=4, dependencies=[[2, 1], [3, 1], [1, 4]], k=2))
# 在第一个学期中，我们可以上课程 2 和课程 3 。然后第二个学期上课程 1 ，第三个学期上课程 4 。
print(Solution().minNumberOfSemesters(n=5, dependencies=[[2, 1], [3, 1], [4, 1], [1, 5]], k=2))
