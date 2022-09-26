"""环形前缀和 断环成链

如果k > 0,将第i个数字用 接下来k个数字之和替换。
如果k < 0,将第i个数字用 之前k个数字之和替换。
如果k == 0,将第i个数字用0替换。
"""

from itertools import accumulate
from typing import List


class Solution:
    def decrypt(self, code: List[int], k: int) -> List[int]:
        """环形数组前缀和"""
        n = len(code)
        nums = code * 2
        preSum = [0] + list(accumulate(nums))
        res = []
        for i in range(len(code)):
            if k > 0:
                res.append(preSum[i + k + 1] - preSum[i + 1])
            elif k < 0:
                res.append(preSum[n + i] - preSum[n + i + k])
            else:
                res.append(0)
        return res


print(Solution().decrypt([5, 7, 1, 4], 3))
# 输入：code = [5,7,1,4], k = 3
# 输出：[12,10,16,13]
# 解释：每个数字都被接下来 3 个数字之和替换。解密后的密码为 [7+1+4, 1+4+5, 4+5+7, 5+7+1]。注意到数组是循环连接的。
# 示例 2：

# 输入：code = [1,2,3,4], k = 0
# 输出：[0,0,0,0]
# 解释：当 k 为 0 时,所有数字都被 0 替换。
# 示例 3：

# 输入：code = [2,4,9,3], k = -2
# 输出：[12,5,6,13]
# 解释：解密后的密码为 [3+9, 2+3, 4+2, 9+4] 。注意到数组是循环连接的。如果 k 是负数,那么和为 之前 的数字。
