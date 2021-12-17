from typing import List
from collections import deque

# k,digit1,digit2
# 找到最小的比k大的，是k倍数的，只有digit1/digit2的数
# 1 <= k <= 1000
# 0 <= digit1 <= 9
# 0 <= digit2 <= 9


# Return the smallest such integer.
# If no such integer exists or the integer exceeds the limit of a signed 32-bit integer (231 - 1),
# return -1.
class Solution:
    def findInteger(self, k: int, digit1: int, digit2: int) -> int:
        if digit1 == 0 and digit2 == 0:
            return -1
        if digit1 > digit2:
            digit1, digit2 = digit2, digit1
        queue = deque([digit1, digit2])

        while queue:
            cur = queue.popleft()
            if cur > (2 ** 31 - 1):
                return -1
            if cur % k == 0 and cur > k:
                return cur
            # 按照顺序加入队列 整个队列都是有序的
            queue.append(cur * 10 + digit1)
            queue.append(cur * 10 + digit2)

        return -1


print(Solution().findInteger(k=3, digit1=4, digit2=2))
print(Solution().findInteger(k=3, digit1=0, digit2=7))
# Output: 24
# Explanation:
# 24 is the first integer larger than 3, a multiple of 3, and comprised of only the digits 4 and/or 2.
