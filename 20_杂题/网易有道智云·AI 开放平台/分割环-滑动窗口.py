# 小易有 n 个数字排成一个环，你能否将它们分成连续的两个部分(即在环上必须连续)，使得两部分的和相等？


from collections import defaultdict
from typing import List


def check(nums: List[int]) -> str:
    sum_ = sum(nums)
    if sum_ & 1:
        return 'NO'
    sum_ >>= 1

    curSum = 0
    left = 0
    for right in range(len(nums)):
        curSum += nums[right]
        if curSum == sum_:
            return 'YES'
        while curSum > sum_:
            curSum -= nums[left]
            left += 1
            if curSum == sum_:
                return 'YES'
    return 'NO'


total_num = int(input())
for i in range(total_num):
    size = input()
    nums = list(map(int, input().split()))
    print(check(nums))

