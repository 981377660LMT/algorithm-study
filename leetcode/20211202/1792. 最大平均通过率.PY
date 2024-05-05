from typing import List
from heapq import heappush, heappop, heapify

# 给你一个整数 extraStudents ，表示额外有 extraStudents 个聪明的学生，
# 他们 一定 能通过任何班级的期末考。你需要给这 extraStudents 个学生每人都安排一个班级，
# 使得 所有 班级的 平均 通过率 最大 。

# 请你返回在安排这 extraStudents 个学生去对应班级后的 最大 平均通过率。

# 总结:
# 坑：不可以找到最小通过率的班级，然后往班级中加学生
# 应该比较增量:谁的通过率增加更多，就越排到堆顶：
class Solution:
    def maxAverageRatio(self, classes: List[List[int]], extraStudents: int) -> float:
        pq = [(good / all - (good + 1) / (all + 1), good, all) for good, all in classes]
        heapify(pq)

        while extraStudents > 0:
            _, good, all = heappop(pq)
            good, all = good + 1, all + 1
            heappush(pq, (good / all - (good + 1) / (all + 1), good, all))
            extraStudents -= 1

        return sum(good / all for _, good, all in pq) / len(classes)


print(Solution().maxAverageRatio(classes=[[1, 2], [3, 5], [2, 2]], extraStudents=2))
# 输出：0.78333
# 解释：你可以将额外的两个学生都安排到第一个班级，平均通过率为 (3/4 + 3/5 + 2/2) / 3 = 0.78333 。
