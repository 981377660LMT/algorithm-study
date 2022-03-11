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

    def getSwapPresum(indexes: List[int]) -> List[int]:
        # 每个字符移动到中间需要的步数，注意左右两边每排好一个位置到中心的距离就要减1，排序，再求前缀和
        mid = len(indexes) // 2
        leftDist, rightDist = (
            [indexes[mid] - indexes[i] for i in range(mid)],
            [indexes[i] - indexes[mid] for i in range(mid, len(indexes))],
        )
        leftDist, rightDist = sorted(leftDist), sorted(rightDist)
        leftDist, rightDist = (
            [num - i for i, num in enumerate(leftDist)],
            [num - i for i, num in enumerate(rightDist)],
        )

        dist = sorted(leftDist + rightDist)
        preSum = list(accumulate(dist))
        return preSum

    indexesMap = defaultdict(list)
    for i, char in enumerate(string):
        indexesMap[char].append(i)

    res = 1
    for char in indexesMap:
        indexes = indexesMap[char]
        if len(indexes) <= res:
            continue

        preSum = getSwapPresum(indexes)

        # 二分查找小于等于k的最大值
        target = bisect_right(preSum, k) - 1
        res = max(res, target + 1)

    return res


string, k = input().split()
k = int(k)
# string, k = 'abcbaa', 2
# string, k = (
#     "zoiumptccefmqdrjhhlgeyljbofwgvwogmvmpzgmoxdrbfdggimzifpfqmrqnrqrlobhluunzhyxrsicdhsrxpsrurqrewvrrcqc",
#     200,
# )
print(maxSame(string, k))
