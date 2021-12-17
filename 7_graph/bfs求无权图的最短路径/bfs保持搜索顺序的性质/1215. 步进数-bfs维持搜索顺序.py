# 给你两个整数，low 和 high，请你找出在 [low, high] 范围内的所有步进数，并返回 排序后 的结果。
# 输入：low = 0, high = 21
# 输出：[0,1,2,3,4,5,6,7,8,9,10,12,21]


# 0 <= low <= high <= 2 * 10^9
# 如果一个整数上的每一位数字与其相邻位上的数字的绝对差都是 1，那么这个数就是一个「步进数」。

from typing import List
from collections import deque

# bfs 层序扩张 省掉了排序的步骤
class Solution:
    def countSteppingNumbers(self, low: int, high: int) -> List[int]:
        res = []
        if low == 0:
            res.append(0)

        queue = deque(list(range(1, 10)))
        while queue:
            cur = queue.popleft()
            if cur > high:
                return res
            if low <= cur:
                res.append(cur)
            last_bit = cur % 10

            # 按照顺序加入队列 整个队列都是有序的
            if last_bit > 0:
                queue.append(cur * 10 + (last_bit - 1))
            if last_bit < 9:
                queue.append(cur * 10 + (last_bit + 1))
        return res


print(Solution().countSteppingNumbers(0, 21))
