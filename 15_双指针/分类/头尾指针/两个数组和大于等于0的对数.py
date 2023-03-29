# 两个数组和大于等于0的对数
# 二分/双指针 这里用头尾双指针解法


from typing import List


def solve(nums1: List[int], nums2: List[int]) -> int:
    nums1, nums2 = sorted(nums1), sorted(nums2)
    res, right = 0, len(nums2) - 1
    for i in range(len(nums1)):
        while right >= 0 and nums1[i] + nums2[right] >= 0:
            right -= 1
        res += len(nums2) - 1 - right
    return res


if __name__ == "__main__":
    import random

    nums1 = [random.randint(-10, 10) for _ in range(10)]
    nums2 = [random.randint(-10, 10) for _ in range(10)]
    assert solve(nums1, nums2) == sum(1 for i in nums1 for j in nums2 if i + j >= 0)
