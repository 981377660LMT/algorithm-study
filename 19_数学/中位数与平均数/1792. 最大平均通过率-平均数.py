# 1792. 最大平均通过率
# https://leetcode.cn/problems/maximum-average-pass-ratio/description/
# 一所学校里有一些班级，每个班级里有一些学生，现在每个班都会进行一场期末考试。给你一个二维数组 classes ，其中 classes[i] = [passi, totali] ，表示你提前知道了第 i 个班级总共有 totali 个学生，其中只有 passi 个学生可以通过考试。
# 给你一个整数 extraStudents ，表示额外有 extraStudents 个聪明的学生，他们 一定 能通过任何班级的期末考。你需要给这 extraStudents 个学生每人都安排一个班级，使得 所有 班级的 平均 通过率 最大 。
#
# !应该比较增量:谁的通过率增加更多，就越排到堆顶

from heapq import heapify, heappop, heappush
from typing import List


class Solution:
    def maxAverageRatio(self, classes: List[List[int]], extraStudents: int) -> float:
        pq = [(good / all - (good + 1) / (all + 1), good, all) for good, all in classes]
        heapify(pq)

        for _ in range(extraStudents):
            _, good, all = heappop(pq)
            good, all = good + 1, all + 1
            heappush(pq, (good / all - (good + 1) / (all + 1), good, all))

        return sum(good / all for _, good, all in pq) / len(classes)
