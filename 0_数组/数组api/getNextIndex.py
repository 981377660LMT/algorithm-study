# getPrevIndex/getNextIndex
# 获取数组中相同元素的前一个和后一个元素的位置


from typing import List, Tuple


def getPrevIndex(arr: List[int]) -> List[int]:
    """
    获取数组中相同元素的前一个元素的位置.不存在则返回-1.
    """

    pool = dict()
    for i, v in enumerate(arr):
        arr[i] = pool.setdefault(v, len(pool))

    n = len(arr)
    nexts, valueNexts = [-1] * n, [-1] * len(pool)
    for i, v in enumerate(arr):
        if valueNexts[v] != -1:
            nexts[i] = valueNexts[v]
        valueNexts[v] = i
    return nexts


def getNextIndex(arr: List[int]) -> List[int]:
    """
    获取数组中相同元素的后一个元素的位置.不存在则返回-1.
    """
    pool = dict()
    for i, v in enumerate(arr):
        arr[i] = pool.setdefault(v, len(pool))

    n = len(arr)
    nexts, valueNexts = [-1] * n, [-1] * len(pool)
    for i in range(n - 1, -1, -1):
        v = arr[i]
        if valueNexts[v] != -1:
            nexts[i] = valueNexts[v]
        valueNexts[v] = i
    return nexts


def getPrevAndNextIndex(arr: List[int]) -> Tuple[List[int], List[int]]:
    """
    获取数组中相同元素的前一个和后一个元素的位置.不存在则返回-1.
    """
    pool = dict()
    for i, v in enumerate(arr):
        arr[i] = pool.setdefault(v, len(pool))

    n = len(arr)
    nexts, prevs, valueNexts, valuePrevs = [-1] * n, [-1] * n, [-1] * len(pool), [-1] * len(pool)
    for i, v in enumerate(arr):
        if valueNexts[v] != -1:
            prevs[i] = valueNexts[v]
        valueNexts[v] = i
    for i in range(n - 1, -1, -1):
        v = arr[i]
        if valuePrevs[v] != -1:
            nexts[i] = valuePrevs[v]
        valuePrevs[v] = i
    return prevs, nexts


if __name__ == "__main__":
    assert getPrevIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == [-1, -1, -1, -1, -1, 0, 1, 2, 3, 4]
    assert getNextIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == [5, 6, 7, 8, 9, -1, -1, -1, -1, -1]
    assert getPrevAndNextIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == (
        [-1, -1, -1, -1, -1, 0, 1, 2, 3, 4],
        [5, 6, 7, 8, 9, -1, -1, -1, -1, -1],
    )
