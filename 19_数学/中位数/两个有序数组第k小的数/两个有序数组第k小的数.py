# 两个有序数组第k小的数/从两个已排序数组中找到第k小的数字
# https://stackoverflow.com/questions/4607945/how-to-find-the-kth-smallest-element-in-the-union-of-two-sorted-arrays/8935157#8935157

from typing import List


def kthSmallestOfTwoSortedArrays(nums1: List[int], nums2: List[int], k: int) -> int:
    """两个有序数组第k(k>=0)小的数."""
    i1, i2 = 0, 0
    j1, j2 = len(nums1), len(nums2)
    while True:
        if i1 == j1:
            return nums2[i2 + k]
        if i2 == j2:
            return nums1[i1 + k]
        mid1, mid2 = (j1 - i1) >> 1, (j2 - i2) >> 1
        if mid1 + mid2 < k:
            if nums1[i1 + mid1] < nums2[i2 + mid2]:
                i1 += mid1 + 1
                k -= mid1 + 1
            else:
                i2 += mid2 + 1
                k -= mid2 + 1
        else:
            if nums1[i1 + mid1] < nums2[i2 + mid2]:
                j2 = i2 + mid2
            else:
                j1 = i1 + mid1


if __name__ == "__main__":

    def brute_force(nums1, nums2, k):
        nums = sorted(nums1 + nums2)
        return nums[k]

    import random

    for _ in range(1000):
        nums1 = sorted([random.randint(0, 100) for _ in range(random.randint(0, 100))])
        nums2 = sorted([random.randint(0, 100) for _ in range(random.randint(0, 100))])
        k = random.randint(0, len(nums1) + len(nums2) - 1)
        a, b = kthSmallestOfTwoSortedArrays(nums1, nums2, k), brute_force(nums1, nums2, k)
        assert a == b, (a, b)

    print("Passed all tests!")
