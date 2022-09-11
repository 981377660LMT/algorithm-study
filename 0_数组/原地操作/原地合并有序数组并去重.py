# 有两个整形数组A和B
# A是有序递增的
# B是有序递减的
# 假设A有额外的空间
# !把B合并到A后需要去重并保证有序递增
# !需要去重


from typing import List


def merge(nums1: List[int], nums2: List[int], n1: int, n2: int) -> List[int]:
    i, j, k = n1 - 1, 0, n1 + n2 - 1
    while j < n2:
        if nums1[i] > nums2[j]:
            nums1[k] = nums1[i]
            i -= 1
        else:
            nums1[k] = nums2[j]
            j += 1
        k -= 1

    # 去重
    slow = 0
    for i in range(n1 + n2):
        if nums1[i] != nums1[slow]:
            slow += 1
            nums1[slow] = nums1[i]

    return nums1[: slow + 1]


print(merge([1, 2, 3, 0, 0, 0], [6, 4, 2], 3, 3))
