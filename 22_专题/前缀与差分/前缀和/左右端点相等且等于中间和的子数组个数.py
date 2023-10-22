# 左右端点相等且等于中间和的好子数组个数

# 给一个数组，定义一个good subarray是最左边的数字等于最右边的数字，
# 还等于中间的数字们的Sum，问有多少个good subarray

# !当前值为num,前缀和为curSum，在哈希表里找值为num前缀和curSum-2*num 的个数


from collections import defaultdict
from typing import List


def countGoodSubarray(nums: List[int]) -> int:
    preSum = defaultdict(dict)
    res, curSum = 0, 0
    for num in nums:
        curSum += num
        mp = preSum[num]
        target = curSum - 2 * num
        if target in mp:
            res += mp[target]
        mp[curSum] = mp.get(curSum, 0) + 1
    return res


if __name__ == "__main__":
    import random

    def bruteForce(nums: List[int]) -> int:
        res = 0
        for i in range(len(nums)):
            for j in range(i + 2, len(nums)):
                if nums[i] == nums[j] and sum(nums[i + 1 : j]) == nums[i]:
                    res += 1
        return res

    nums = [random.randint(-100, 100) for _ in range(1000)]
    assert countGoodSubarray(nums) == bruteForce(nums)
    print("ok")
