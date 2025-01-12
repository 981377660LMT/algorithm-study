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


if __name__ == "__main__":
    from random import randint
    from itertools import combinations

    def count132BruteForce(nums: List[int]) -> int:
        res = 0
        for i, j, k in combinations(range(len(nums)), 3):
            res += nums[i] < nums[k] < nums[j]
        return res

    # random check
    for _ in range(10000):
        n = randint(1, 10)
        nums = [randint(-10, 10) for _ in range(n)]

        if count132(nums) != count132BruteForce(nums):
            print(nums)
            break

    print("Passed")
