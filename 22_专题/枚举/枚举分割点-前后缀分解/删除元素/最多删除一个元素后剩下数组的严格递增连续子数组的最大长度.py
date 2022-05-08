from typing import List


def maxLenAfterRemove(nums: List[int]) -> int:
    """给出一个数组，最多删除一个元素，求剩下数组的严格递增连续子数组的最大长度。
    
    n<=1e6
    """
    n = len(nums)
    pre = [1] * n
    suf = [1] * n
    for i in range(1, n):
        if nums[i] > nums[i - 1]:
            pre[i] = pre[i - 1] + 1

    for i in range(n - 2, -1, -1):
        if nums[i] < nums[i + 1]:
            suf[i] = suf[i + 1] + 1

    res = 1

    # 不删除
    curMax = 1
    for i in range(1, n):
        if nums[i] > nums[i - 1]:
            curMax += 1
        else:
            curMax = 1
        res = max(res, curMax)

    # 枚举删除每个元素
    for i in range(1, n - 1):
        if nums[i - 1] < nums[i + 1]:
            res = max(res, pre[i - 1] + suf[i + 1])

    return res


print(maxLenAfterRemove([1, 2, 5, 3, 4]))  # 4
print(maxLenAfterRemove([1, 2]))  # 2
