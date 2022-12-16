# 给你一串数，每个数都只能做将高度减1和不改变两种操作，问这串数是否可以变成不减序列

from typing import List


def buildStairs(nums: List[int]) -> bool:
    """倒序遍历"""
    for i in range(len(nums) - 2, -1, -1):
        if nums[i] > nums[i + 1]:
            nums[i] -= 1
        if nums[i] > nums[i + 1]:
            return False
    return True


n = int(input())
nums = list(map(int, input().split()))
print("Yes" if buildStairs(nums) else "No")
