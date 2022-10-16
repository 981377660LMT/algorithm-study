# 768. 最多能完成排序的块 II
#  arr是一个可能包含重复元素的整数数组，
#  我们将这个数组分割成几个“块”，并将这些块分别进行排序。
#  !之后再连接起来，使得连接的结果和按升序排序后的原数组相同。
#  我们最多能将数组分成多少块？
from itertools import accumulate
from typing import List

INF = int(1e18)


class Solution:
    def maxChunksToSorted(self, arr: List[int]) -> int:
        """`前缀的最大值小于等于后缀的最小值`的位置的数量。"""
        sufMin = ([INF] + list(accumulate(arr[::-1], min)))[::-1]
        preMax = [-INF] + list(accumulate(arr, max))
        return sum(preMax[i] <= sufMin[i] for i in range(1, len(arr) + 1))

    def maxChunksToSorted2(self, arr: List[int]) -> int:
        """排序+前缀和比较几个点相等"""
        preSum1 = list(accumulate(arr))
        preSum2 = list(accumulate(sorted(arr)))
        return sum(preSum1[i] == preSum2[i] for i in range(len(arr)))


print(Solution().maxChunksToSorted([5, 4, 3, 2, 1]))
# 输出: 1
# 解释:
# 将数组分成2块或者更多块，都无法得到所需的结果。
# 例如，分成 [5, 4], [3, 2, 1] 的结果是 [4, 5, 1, 2, 3]，这不是有序的数组。

print(Solution().maxChunksToSorted([2, 1, 3, 4, 4]))
# 输出: 4
# 解释:
# 我们可以把它分成两块，例如 [2, 1], [3, 4, 4]。
# 然而，分成 [2, 1], [3], [4], [4] 可以得到最多的块数。
