from typing import List
from bisect import bisect_left

# difficulty[i] 表示第 i 个工作的难度，profit[i] 表示第 i 个工作的收益。
# worker[i] 是第 i 个工人的能力，即该工人只能完成难度小于等于 worker[i] 的工作。
# !每一个工人都最多只能安排一个工作，但是一个工作可以完成多次。
# 我们能得到的最大收益是多少？(调整打怪策略，与什么样的怪兽战斗获得金币最多)
# 1 <= n, m <= 104

# 不能二分查找：题目没有说难度高的工作收益也越大

# 总结：
# !排序+遍历加指针记录
# !对每个工人找到最大收益


class Solution:
    def maxProfitAssignment(
        self, difficulty: List[int], profit: List[int], worker: List[int]
    ) -> int:
        worker.sort()
        jobs = sorted(zip(difficulty, profit), key=lambda x: x[0])
        res, ji, maxPro = 0, 0, 0
        for w in worker:
            while ji < len(jobs) and jobs[ji][0] <= w:
                maxPro = max(maxPro, jobs[ji][1])
                ji += 1
            res += maxPro
        return res


print(
    Solution().maxProfitAssignment(
        difficulty=[2, 4, 6, 8, 10], profit=[10, 20, 30, 40, 50], worker=[4, 5, 6, 7]
    )
)
# 输出: 100
# 解释: 工人被分配的工作难度是 [4,4,6,6] ，分别获得 [20,20,30,30] 的收益
