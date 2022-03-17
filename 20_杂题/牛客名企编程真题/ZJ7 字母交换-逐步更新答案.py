# 字符串S由小写字母构成，长度为n。定义一种操作，
# 每次都可以挑选字符串中任意的两个相邻字母进行交换。
# 询问在至多交换m次之后，字符串中最多有多少个连续的位置上的字母相同？

# 第一行为一个字符串S与一个非负整数m。(1 <= |S| <= 1e3, 1 <= m <= 1e6)

# 对于每一个字符，如果这个key值下数据长度短于目前找到的最大长度，就直接跳过。
# 否则遍历（目前最大长度，该key值数据长度），寻找是否有符合条件的更优解。
# 移动数据的时候，所有数据往中间数据靠是最优解。

from typing import List
from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate


def maxSame(string: str, k: int) -> int:
    """在至多交换m次之后，字符串中最多有多少个连续的位置上的字母相同"""

    def minMoves(indexes: List[int], target: int) -> int:
        """得到连续 target 个 相同字符 的最少相邻交换次数"""
        indexes = [num - i for i, num in enumerate(indexes)]
        preSum = [0] + list(accumulate(indexes))

        res = int(1e20)
        # 把ones里的哪k个数移动到一起  left+k-1<len(ones)
        for left in range(len(indexes) + 1 - target):
            right = left + target - 1
            mid = (left + right) >> 1
            # mid左右两边的和
            leftSum = indexes[mid] * (mid - left) - (preSum[mid] - preSum[left])
            rightSum = preSum[right + 1] - preSum[mid + 1] - indexes[mid] * (right - mid)
            res = min(res, leftSum + rightSum)

        return res

    indexesMap = defaultdict(list)
    for i, char in enumerate(string):
        indexesMap[char].append(i)

    res = 1
    for char in indexesMap:
        indexes = indexesMap[char]
        if len(indexes) <= res:
            continue

        left, right = 1, len(indexes)
        while left <= right:
            mid = (left + right) >> 1
            # 得到连续mid个相同字符的最小交换次数
            if minMoves(indexes, mid) <= k:
                left = mid + 1
            else:
                right = mid - 1

        res = max(res, right)

    return res


string, k = input().split()
k = int(k)
# string, k = 'abcbaa', 2
# string, k = (
#     "zoiumptccefmqdrjhhlgeyljbofwgvwogmvmpzgmoxdrbfdggimzifpfqmrqnrqrlobhluunzhyxrsicdhsrxpsrurqrewvrrcqc",
#     200,
# )
print(maxSame(string, k))
