#     1
#     1
#     1    1
#     1    1
# 1   1    1
# 1   1    1
# 是否存在132模式的子序列
# 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，
# 并同时满足：i < j < k 和 nums[i] < nums[k] < nums[j] 。
# !n<=2e5


from typing import List
from sortedcontainers import SortedList

INF = int(1e20)


def count132(nums: List[int]) -> int:
    """
    统计132模式的个数 O(nlogn)
    nums[i] < nums[k] < nums[j]

    先计算出123模式、132模式、122模式的个数之和,再减去123模式和122模式的个数
    """
    n = len(nums)
    sum_ = 0
    right = SortedList(nums)
    for i in range(n):
        right.remove(nums[i])
        count = len(right) - right.bisect_right(nums[i])
        sum_ += count * (count - 1) // 2
    return sum_ - count123122(nums)


def count123122(nums: List[int]) -> int:
    """统计123模式与122模式的个数 O(nlogn)

    nums[i]<nums[j]<=nums[k]
    """
    res = 0
    leftSmaller = SortedList()
    rightBigger = SortedList(nums)
    for i in range(len(nums)):
        rightBigger.remove(nums[i])
        count1 = leftSmaller.bisect_left(nums[i])  # 左侧严格小于nums[i]的个数
        count2 = len(rightBigger) - rightBigger.bisect_left(nums[i])  # !右侧大于等于nums[i]的个数
        res += count1 * count2
        leftSmaller.add(nums[i])
    return res


def count123(nums: List[int]) -> int:
    """统计123模式的个数 O(nlogn)

    nums[i]<nums[j]<nums[k]
    """
    res = 0
    leftSmaller = SortedList()
    rightBigger = SortedList(nums)
    for i in range(len(nums)):
        rightBigger.remove(nums[i])
        count1 = leftSmaller.bisect_left(nums[i])  # 左侧严格小于nums[i]的个数
        count2 = len(rightBigger) - rightBigger.bisect_right(nums[i])  # 右侧严格大于nums[i]的个数
        res += count1 * count2
        leftSmaller.add(nums[i])
    return res


def count1324(nums: List[int]) -> int:
    """统计1324模式的个数 O(n^2)

    nums[i] < nums[k] < nums[j] < nums[l]

    枚举第二个（"i2"）和第四个（"i4"）元素，然后统计第一个（"i1"）和第三个（"i3"）元素的个数
    然后对于每个 "i2" "i4" 对，统计以 "i2" 为 "中间" 的、且在 "i4" 之前的 132 模式数量。
    """
    n = len(nums)
    res = 0
    counter = [0] * n  # 以 i2 为 "中间" 的 "132" 模式的数量
    for i4 in range(n):
        bigger = 0
        for i2 in range(i4):
            if nums[i4] > nums[i2]:
                res += counter[i2]
                bigger += 1
            else:
                counter[i2] += bigger
    return res


if __name__ == "__main__":
    from random import randint
    from itertools import combinations

    def count132BruteForce(nums: List[int]) -> int:
        res = 0
        for i, j, k in combinations(range(len(nums)), 3):
            res += nums[i] < nums[k] < nums[j]
        return res

    # random check
    while True:
        n = randint(1, 10)
        nums = [randint(-10, 10) for _ in range(n)]
        print(1)
        if count132(nums) != count132BruteForce(nums):
            print(nums)
            break
