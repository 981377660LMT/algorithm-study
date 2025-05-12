# 这个问题如何优化
# for v1 in dist1:
#     for v2 in dist2:
#         res += max(d2, v1 + v2 + 1)


from bisect import bisect_right
from itertools import accumulate
from typing import List


def convolutionChmax1(arr1: List[int], arr2: List[int], chMax: int) -> int:
    sorted2 = sorted(arr2)
    presum2 = list(accumulate(sorted2, initial=0))
    res = 0
    for v1 in arr1:
        threshold = chMax - v1
        pos = bisect_right(sorted2, threshold)
        res += chMax * pos
        if pos < len(sorted2):
            remainingSum = presum2[-1] - presum2[pos]
            remainingCount = len(sorted2) - pos
            res += remainingSum + v1 * remainingCount
    return res


def convolutionChmax2(arr1: List[int], arr2: List[int], chMax: int) -> int:
    arr1, arr2 = sorted(arr1, reverse=True), sorted(arr2)
    presum2 = list(accumulate(arr2, initial=0))
    res = 0
    pos = 0
    for v1 in arr1:
        threshold = chMax - v1
        while pos < len(arr2) and arr2[pos] <= threshold:
            pos += 1
        res += chMax * pos
        if pos < len(arr2):
            remainingSum = presum2[-1] - presum2[pos]
            remainingCount = len(arr2) - pos
            res += remainingSum + v1 * remainingCount
    return res


if __name__ == "__main__":

    def checkWithBruteForce():
        from itertools import product
        from random import randint

        for _ in range(100):
            n1 = randint(1, 10)
            n2 = randint(1, 10)
            arr1 = [randint(1, 10) for _ in range(n1)]
            arr2 = [randint(1, 10) for _ in range(n2)]
            chMax = randint(0, 20)
            bruteForce = sum(max(chMax, v1 + v2) for v1, v2 in product(arr1, arr2))
            res1 = convolutionChmax1(arr1, arr2, chMax)
            res2 = convolutionChmax2(arr1, arr2, chMax)
            assert (
                res1 == res2 == bruteForce
            ), f"res1: {res1}, res2: {res2}, bruteForce: {bruteForce}, arr1: {arr1}, arr2: {arr2}, chMax: {chMax}"
        print("Passed all test cases.")

    checkWithBruteForce()
