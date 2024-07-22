# 两个数组的交集/并集/补集/差集 (交并补差)
# api:
#   intersection(a: List[int], b: List[int]) -> List[int]
#   union(a: List[int], b: List[int]) -> List[int]
#   difference(a: List[int], b: List[int]) -> List[int]
#   symmetricDifference(a: List[int], b: List[int]) -> List[int]
#   splitDifferenceAndIntersection(a: List[int], b: List[int]) -> Tuple[List[int], List[int], List[int]]
#   isSubset(a: List[int], b: List[int]) -> bool
#   isDisjoint(a: List[int], b: List[int]) -> bool


from typing import List, Tuple


def intersection(a: List[int], b: List[int]) -> List[int]:
    """两个有序数组的交集."""
    i, n = 0, len(a)
    j, m = 0, len(b)
    res = []
    while True:
        if i == n or j == m:
            return res
        x, y = a[i], b[j]
        if x < y:
            i += 1
        elif x > y:
            j += 1
        else:
            res.append(a[i])
            i += 1
            j += 1


def union(a: List[int], b: List[int]) -> List[int]:
    """两个有序数组的并集(合并两个有序数组)."""
    return merge(a, b)


def merge(a: List[int], b: List[int]) -> List[int]:
    i, n = 0, len(a)
    j, m = 0, len(b)
    res = []
    while True:
        if i == n:
            res += b[j:]
            return res
        if j == m:
            res += a[i:]
            return res
        x, y = a[i], b[j]
        if x < y:
            res.append(a[i])
            i += 1
        elif x > y:
            res.append(b[j])
            j += 1
        else:
            res.append(a[i])
            i += 1
            j += 1


def difference(a: List[int], b: List[int]) -> List[int]:
    """两个有序数组的差集 a-b."""
    i, n = 0, len(a)
    j, m = 0, len(b)
    res = []
    while True:
        if i == n:
            return res
        if j == m:
            res += a[i:]
            return res
        x, y = a[i], b[j]
        if x < y:
            res.append(a[i])
            i += 1
        elif x > y:
            j += 1
        else:
            i += 1
            j += 1


def symmetricDifference(a: List[int], b: List[int]) -> List[int]:
    """两个有序数组的对称差集 a▲b."""
    i, n = 0, len(a)
    j, m = 0, len(b)
    res = []
    while True:
        if i == n:
            res += b[j:]
            return res
        if j == m:
            res += a[i:]
            return res
        x, y = a[i], b[j]
        if x < y:
            res.append(a[i])
            i += 1
        elif x > y:
            res.append(b[j])
            j += 1
        else:
            i += 1
            j += 1


def splitDifferenceAndIntersection(
    a: List[int], b: List[int]
) -> Tuple[List[int], List[int], List[int]]:
    """求差集 A-B, B-A 和交集 A∩B."""
    differenceA, differenceB, intersection = [], [], []
    i, n = 0, len(a)
    j, m = 0, len(b)
    while True:
        if i == n:
            differenceB += b[j:]
            return differenceA, differenceB, intersection
        if j == m:
            differenceA += a[i:]
            return differenceA, differenceB, intersection
        x, y = a[i], b[j]
        if x < y:
            differenceA.append(x)
            i += 1
        elif x > y:
            differenceB.append(y)
            j += 1
        else:
            intersection.append(x)
            i += 1
            j += 1


def isSubset(a: List[int], b: List[int]) -> bool:
    """a 是否为 b 的子集（相当于 differenceA 为空）."""
    i, n = 0, len(a)
    j, m = 0, len(b)
    while True:
        if i == n:
            return True
        if j == m:
            return False
        x, y = a[i], b[j]
        if x < y:
            return False
        elif x > y:
            j += 1
        else:
            i += 1
            j += 1


def isDisjoint(a: List[int], b: List[int]) -> bool:
    """是否为不相交集合（相当于 intersection 为空."""
    i, n = 0, len(a)
    j, m = 0, len(b)
    while True:
        if i == n or j == m:
            return True
        x, y = a[i], b[j]
        if x < y:
            i += 1
        elif x > y:
            j += 1
        else:
            return False


if __name__ == "__main__":
    assert intersection([1, 2, 3, 4], [3, 4, 5, 6]) == [3, 4]
    assert union([1, 2, 3, 4], [3, 4, 5, 6]) == [1, 2, 3, 4, 5, 6]
    assert difference([3, 4, 5, 6], [1, 2, 3, 4]) == [5, 6]
    assert symmetricDifference([1, 2, 3, 4], [3, 4, 5, 6]) == [1, 2, 5, 6]
    assert splitDifferenceAndIntersection([1, 2, 3, 4], [3, 4, 5, 6]) == ([1, 2], [5, 6], [3, 4])
    assert isSubset([1, 2, 3], [1, 2, 3, 4, 5])
    assert not isDisjoint([1, 2, 3, 4], [4, 5, 6])
    print("PASSED")
