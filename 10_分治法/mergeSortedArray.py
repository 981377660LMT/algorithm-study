from typing import List


def mergeTwoSortedArray(arr1: List[int], arr2: List[int]) -> List[int]:
    """合并两个有序数组."""
    n1, n2 = len(arr1), len(arr2)
    if n1 == 0:
        return arr2
    if n2 == 0:
        return arr1
    res = [0] * (n1 + n2)
    i, j, k = 0, 0, 0
    while i < n1 and j < n2:
        if arr1[i] < arr2[j]:
            res[k] = arr1[i]
            i += 1
        else:
            res[k] = arr2[j]
            j += 1
        k += 1
    while i < n1:
        res[k] = arr1[i]
        i += 1
        k += 1
    while j < n2:
        res[k] = arr2[j]
        j += 1
        k += 1
    return res


def mergeThreeSortedArray(arr1: List[int], arr2: List[int], arr3: List[int]) -> List[int]:
    """合并三个有序数组."""
    n1, n2, n3 = len(arr1), len(arr2), len(arr3)
    if n1 == 0:
        return mergeTwoSortedArray(arr2, arr3)
    if n2 == 0:
        return mergeTwoSortedArray(arr1, arr3)
    if n3 == 0:
        return mergeTwoSortedArray(arr1, arr2)
    res = [0] * (n1 + n2 + n3)
    i1, i2, i3, k = 0, 0, 0, 0
    while i1 < n1 and i2 < n2 and i3 < n3:
        if arr1[i1] < arr2[i2]:
            if arr1[i1] < arr3[i3]:
                res[k] = arr1[i1]
                i1 += 1
            else:
                res[k] = arr3[i3]
                i3 += 1
        else:
            if arr2[i2] < arr3[i3]:
                res[k] = arr2[i2]
                i2 += 1
            else:
                res[k] = arr3[i3]
                i3 += 1
        k += 1
    while i1 < n1 and i2 < n2:
        if arr1[i1] < arr2[i2]:
            res[k] = arr1[i1]
            i1 += 1
        else:
            res[k] = arr2[i2]
            i2 += 1
        k += 1
    while i1 < n1 and i3 < n3:
        if arr1[i1] < arr3[i3]:
            res[k] = arr1[i1]
            i1 += 1
        else:
            res[k] = arr3[i3]
            i3 += 1
        k += 1
    while i2 < n2 and i3 < n3:
        if arr2[i2] < arr3[i3]:
            res[k] = arr2[i2]
            i2 += 1
        else:
            res[k] = arr3[i3]
            i3 += 1
        k += 1
    while i1 < n1:
        res[k] = arr1[i1]
        i1 += 1
        k += 1
    while i2 < n2:
        res[k] = arr2[i2]
        i2 += 1
        k += 1
    while i3 < n3:
        res[k] = arr3[i3]
        i3 += 1
        k += 1
    return res


def mergeKSortedArray(arrays: List[List[int]]) -> List[int]:
    """合并k个有序数组."""
    n = len(arrays)
    if n == 0:
        return []
    if n == 1:
        return arrays[0]
    if n == 2:
        return mergeTwoSortedArray(arrays[0], arrays[1])

    def merge(start: int, end: int) -> List[int]:
        """
        合并[start,end)区间的数组.
        TODO:使用`mergeThreeSortedArray`加速.
        """
        if start >= end:
            return []
        if end - start == 1:
            return arrays[start]
        mid = (start + end) >> 1
        return mergeTwoSortedArray(merge(start, mid), merge(mid, end))

    return merge(0, len(arrays))


if __name__ == "__main__":
    arr1 = [1, 3, 5, 7, 9]
    arr2 = [2, 4, 6, 8, 10]
    arr3 = [0, 11, 12, 13, 14]
    print(mergeThreeSortedArray(arr1, arr2, arr3))
