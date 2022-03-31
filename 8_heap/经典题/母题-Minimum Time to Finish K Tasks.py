# 完成k个任务的最短时间
# tasks为nx3的矩阵

# k ≤ n ≤ 1,000

# 如果只有1列,排序选择前k行即可
# 如果只有2列,第一列排序后,第二列用容量k的大顶堆维护即可


# 推广到任意列的解法
from heapq import heapify, heappushpop


class Solution:
    def solve(self, tasks, k):
        """
        k任务问题
        从矩阵中选择k行组成一个子矩阵
        求这个子矩阵`每列最大值的和`的最小值
        """
        if k == 0:
            return 0

        row = len(tasks)
        col = len(tasks[0])
        tasks.sort(key=lambda x: x[0])
        if col == 1:
            return tasks[k - 1][0]

        res = int(1e20)
        if col == 2:
            pq = [-tasks[i][1] for i in range(k)]
            heapify(pq)
            for i in range(k - 1, row):
                res = min(res, tasks[i][0] - pq[0])
                if i + 1 < row:
                    heappushpop(pq, -tasks[i + 1][1])
        else:
            pq = [tasks[i][1:] for i in range(k)]
            for i in range(k - 1, row):
                res = min(res, tasks[i][0] + self.solve(pq, k))
                if i + 1 < row:
                    pq.append(tasks[i + 1][1:])

        return res


print(Solution().solve(tasks=[[1, 2, 2], [3, 4, 1], [3, 1, 2]], k=2))
# We pick the first row and the last row. And the total sum becomes

# S = [[1,2,2],[3,1,2]]

# max(S[0][0], S[1][0]) = 3
# max(S[0][1], S[1][1]) = 2
# max(S[0][2], S[1][2]) = 2
