from typing import List
from heapq import heapify, heappop, heappush

# 一开始，你有一个数组 A ，它的所有元素均为 1 ，你可以执行以下操作：
# 令 x 为你数组里所有元素的和
# 选择满足 0 <= i < target.size 的任意下标 i ，并让 A 数组里下标为 i 处的值为 x 。
# 你可以重复该过程任意次

# 如果能从 A 开始构造出目标数组 target ，请你返回 True ，否则返回 False 。
# 1 <= target.length <= 5 * 10^4


# 此题关键:`这一步加上去的那个数一定是最大的数`
# 注意几个特殊情况
# [2]
# other_sum 会为0
# [2,1]
# other_sum 会为1
# [10, 1, 1]
# t % other_sum 会取到0
class Solution:
    def isPossible(self, target: List[int]) -> bool:
        total = sum(target)
        pq = [-v for v in target]
        heapify(pq)

        while pq:
            cur = -heappop(pq)
            other_sum = total - cur
            if cur == 1 or other_sum == 1:
                return True
            if cur < other_sum or other_sum == 0 or cur % other_sum == 0:
                return False
            cur %= other_sum  # 细节  [2, 90000002] 这种情况可以加速判断
            total = other_sum + cur
            heappush(pq, -cur)

        return False


print(Solution().isPossible(target=[9, 3, 5]))
# 输出：true
# 解释：从 [1, 1, 1] 开始
# [1, 1, 1], 和为 3 ，选择下标 1
# [1, 3, 1], 和为 5， 选择下标 2
# [1, 3, 5], 和为 9， 选择下标 0
# [9, 3, 5] 完成
print(Solution().isPossible(target=[1, 1, 1, 2]))
print(Solution().isPossible(target=[1, 1, 2]))
print(Solution().isPossible(target=[8, 5]))
