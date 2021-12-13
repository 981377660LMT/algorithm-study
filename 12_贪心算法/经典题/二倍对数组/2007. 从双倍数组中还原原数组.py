from typing import List
from collections import deque

# 954. 二倍数对数组.py
class Solution:
    def findOriginalArray(self, changed: List[int]) -> List[int]:
        n = len(changed)
        if n & 1:
            return []

        changed = sorted(changed)
        res, queue = [], deque()

        # queue里存储一倍,等待num来配对
        for num in changed:
            if queue and num == queue[0] * 2:
                res.append(queue.popleft())
            else:
                queue.append(num)

        return res if len(res) == n // 2 else []


print(Solution().findOriginalArray(changed=[1, 3, 4, 2, 6, 8]))
# 输出：[1,3,4]
# 解释：一个可能的 original 数组为 [1,3,4] :
# - 将 1 乘以 2 ，得到 1 * 2 = 2 。
# - 将 3 乘以 2 ，得到 3 * 2 = 6 。
# - 将 4 乘以 2 ，得到 4 * 2 = 8 。
# 其他可能的原数组方案为 [4,3,1] 或者 [3,1,4] 。

