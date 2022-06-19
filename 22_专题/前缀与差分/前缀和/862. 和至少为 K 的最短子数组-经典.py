from typing import Deque, List, Tuple
from collections import deque

# 1 <= A.length <= 50000
# 带限制的子数组/序列问题
# 1.如果数组中的数据均为非负数的话，那么就对应常规的子数组和问题，可以使用滑动窗口来解决
# 209. 长度最小的子数组
# 但是添加了负数之后，窗口的滑动便丢失了单向性，因此无法使用滑动窗口解决。
Index = int
PreSum = int


class Solution:
    def shortestSubarray(self, nums: List[int], k: int) -> int:
        # 存储前缀和，单增；如果加入的前缀和减去队首的前缀和>=k 那么队首就找到了以他开始的最短的子数组，队首就可以退位了
        # 如果加入的前缀和小于队尾的前缀和 直接删除队尾 因为队尾找到符合题意的子数组还得比后面多带个负数 肯定不是最短的
        queue: Deque[Tuple[Index, PreSum]] = deque([(-1, 0)])
        res = int(1e20)
        cur = 0

        for i, num in enumerate(nums):
            cur += num
            while queue and cur - queue[0][1] >= k:
                preI, _ = queue.popleft()
                res = min(res, i - preI)
            while queue and cur <= queue[-1][1]:
                queue.pop()
            queue.append((i, cur))

        return res if res != int(1e20) else -1

