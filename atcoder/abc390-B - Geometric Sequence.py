# B - Geometric Sequence
# https://atcoder.jp/contests/abc390/tasks/abc390_b

from typing import List


def isGeomerticSequence(arr: List[int]) -> bool:
    """判断是否是等比数列."""
    n = len(arr)
    for i in range(n - 2):
        if arr[i] * arr[i + 2] != arr[i + 1] * arr[i + 1]:
            return False
    return True


if __name__ == "__main__":
    n = int(input())
    arr = list(map(int, input().split()))
    print("Yes" if isGeomerticSequence(arr) else "No")
