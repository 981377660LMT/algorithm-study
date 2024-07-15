from typing import List, Tuple


def mergeTwoArrayWithMinCost(arr1: List[int], arr2: List[int]) -> Tuple[int, List[Tuple[int, int]]]:
    """
    合并两个数组，最小化总代价.
    第i个数的代价为:`(这个数之前"不同"类型的数的个数+1)*这个数`.
    返回最小代价和合并的结果，结果每一项形如(type, index).
    """
    order1 = sorted(list(range(len(arr1))), key=lambda i: -arr1[i])
    order2 = sorted(list(range(len(arr2))), key=lambda i: -arr2[i])
    resCost, resArray = 0, []
    count1, count2 = 1, 1
    i, j = 0, 0
    while i < len(arr1) and j < len(arr2):
        ptr1, ptr2 = order1[i], order2[j]
        if arr1[ptr1] > arr2[ptr2]:
            resCost += arr1[ptr1] * count2
            resArray.append((1, ptr1))
            i += 1
            count1 += 1
        else:
            resCost += arr2[ptr2] * count1
            resArray.append((2, ptr2))
            j += 1
            count2 += 1
    while i < len(arr1):
        ptr1 = order1[i]
        resCost += arr1[ptr1] * count2
        resArray.append((1, ptr1))
        i += 1
    while j < len(arr2):
        ptr2 = order2[j]
        resCost += arr2[ptr2] * count1
        resArray.append((2, ptr2))
        j += 1
    return resCost, resArray
