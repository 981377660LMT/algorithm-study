from typing import List
from sortedcontainers import SortedList

# 每个工人只能完成 一个 任务，且力量值需要 大于等于 该任务的力量要求值（即 workers[j] >= tasks[i] ）。
# 你还有 pills 个神奇药丸，可以给 一个工人的力量值 增加 strength 。你可以决定给哪些工人使用药丸，但每个工人 最多 只能使用 一片 药丸。
# 请你返回 最多 有多少个任务可以被完成。

# 二分查找完成的任务数 + 固定的任务数下贪心选择厉害的工人+药丸分配给离tasks差最小的人
class Solution:
    def maxTaskAssign(self, tasks: List[int], workers: List[int], pills: int, strength: int) -> int:
        tasks.sort()
        workers.sort()

        def check(mid: int) -> bool:
            remain = pills
            #  工人的有序集合
            sw = SortedList(workers[-mid:])
            # 从大到小枚举每一个任务
            for t in reversed(tasks[:mid]):
                # 如果有序集合中最大的元素大于等于最大任务
                if sw[-1] >= t:
                    sw.pop()
                else:
                    if remain == 0:
                        return False
                    # 在有序集合中找出`最小的`大于等于 t−strength 的元素并删除
                    cand = sw.bisect_left(t - strength)
                    if cand == len(sw):
                        return False
                    remain -= 1
                    sw.pop(cand)

            return True

        left, right = 0, min(len(tasks), len(workers))
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maxTaskAssign(tasks=[3, 2, 1], workers=[0, 3, 3], pills=1, strength=1))
# 输出：3
# 解释：
# 我们可以按照如下方案安排药丸：
# - 给 0 号工人药丸。
# - 0 号工人完成任务 2（0 + 1 >= 1）
# - 1 号工人完成任务 1（3 >= 2）
# - 2 号工人完成任务 0（3 >= 3）
