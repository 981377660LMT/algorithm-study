# 2071. 你可以安排的最多任务数目
# 每个工人只能完成 一个 任务，且力量值需要 大于等于 该任务的力量要求值（即 workers[j] >= tasks[i] ）。
# 你还有 pills 个神奇药丸，可以给 一个工人的力量值 增加 strength 。你可以决定给哪些工人使用药丸，但每个工人 最多 只能使用 一片 药丸。
# 请你返回 最多 有多少个任务可以被完成。
#
# 二分查找完成的任务数 + 固定的任务数下贪心选择厉害的工人+药丸分配给离tasks差最小的人
# 双序列配对
# 二分+双端队列
# !能不吃药就不吃药，能吃药就吃最难的任务.


from collections import deque
from typing import List


class Solution:
    def maxTaskAssign(self, tasks: List[int], workers: List[int], pills: int, strength: int) -> int:
        def check(mid: int) -> bool:
            """
            选择mid个任务和mid个工人，判断是否可以完成配对.
            选择最强的 k 名工人，去完成最简单的 k 个任务.
            枚举工人，计算他完成哪个任务.

            如果 w 不吃药，能完成目前剩余任务中最简单的任务，那么就完成最简单的任务.
            如果 w 必须吃药，贪心地，让 w `完成他能完成的最难的任务`，充分利用这颗药的效果.
            !能不吃药就不吃药，能吃药就吃最难的任务.
            """
            ti, p = 0, pills
            validTask = deque()
            for w in workers[-mid:]:
                # 吃药后可以完成的任务
                while ti < mid and tasks[ti] <= w + strength:
                    validTask.append(tasks[ti])
                    ti += 1
                if not validTask:
                    return False
                # 不吃药，完成最简单的任务
                if w >= validTask[0]:
                    validTask.popleft()
                    continue
                # 必须吃药
                if p == 0:
                    return False
                p -= 1
                # 吃药后可以完成的最难的任务
                validTask.pop()
            return True

        tasks.sort()
        workers.sort()

        left, right = 1, min(len(tasks), len(workers))
        while left <= right:
            mid = (left + right) // 2
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
