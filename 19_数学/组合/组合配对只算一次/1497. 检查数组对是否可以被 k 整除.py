from typing import List
from collections import Counter

# 现在需要把数组恰好分成 n / 2 对，以使每对数字的和都能够被 k 整除。
# 1 <= k <= 10^5


# 1. 模1和模k-1配对 ，模2和模k-2配对...
# 2. 模0的个数需要为偶数
class Solution:
    def canArrange(self, arr: List[int], k: int) -> bool:
        modCounter = Counter((num % k for num in arr))
        return modCounter[0] & 1 == 0 and all(
            modCounter[k - mod] == count for mod, count in modCounter.items() if mod > 0
        )


print(Solution().canArrange(arr=[1, 2, 3, 4, 5, 10, 6, 7, 8, 9], k=5))
# 输出：true
# 解释：划分后的数字对为 (1,9),(2,8),(3,7),(4,6) 以及 (5,10)

