from typing import List


def validMountainArray(
    arr: List[int], leftStrict=True, rightStrict=True, allowEmptySide=False
) -> bool:
    """有效的山脉数组."""
    n = len(arr)
    ptr = 0

    if leftStrict:
        while ptr + 1 < n and arr[ptr] < arr[ptr + 1]:
            ptr += 1
    else:
        while ptr + 1 < n and arr[ptr] <= arr[ptr + 1]:
            ptr += 1

    if not allowEmptySide and (ptr == 0 or ptr == n - 1):
        return False

    if rightStrict:
        while ptr + 1 < n and arr[ptr] > arr[ptr + 1]:
            ptr += 1
    else:
        while ptr + 1 < n and arr[ptr] >= arr[ptr + 1]:
            ptr += 1

    return ptr == n - 1


if __name__ == "__main__":
    assert not validMountainArray([2, 1])
    assert not validMountainArray([3, 5, 5])
    assert validMountainArray([0, 3, 2, 1])
    assert not validMountainArray([0, 1, 2, 3, 4, 5, 6, 7, 8, 9])
    assert validMountainArray([0, 1, 2, 3, 4, 5, 6, 7, 8, 9], allowEmptySide=True)
