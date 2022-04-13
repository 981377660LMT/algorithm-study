from typing import List
from bisect import bisect_left, bisect_right


# LIS模板
def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    res = [nums[0]]
    for i in range(1, n):
        pos = bisect_left(res, nums[i]) if isStrict else bisect_right(res, nums[i])
        if pos >= len(res):
            res.append(nums[i])
        else:
            res[pos] = nums[i]

    return len(res)


if __name__ == '__main__':
    assert LIS([10, 9, 2, 5, 3, 7, 101, 18]) == 4

