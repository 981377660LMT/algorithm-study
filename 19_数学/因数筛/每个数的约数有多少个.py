from typing import List


def countFactor(nums: List[int]) -> List[int]:
    """对于每个数，原数组中有多少个他的约数."""
    upper = max(nums) + 1
    c1, c2 = [0] * upper, [0] * upper
    for v in nums:
        c1[v] += 1
    for f in range(1, upper):
        for m in range(f, upper, f):
            c2[m] += c1[f]
    return c2


if __name__ == "__main__":
    # D - Not Divisible
    # https://atcoder.jp/contests/abc170/tasks/abc170_d
    N = int(input())
    A = list(map(int, input().split()))
    counter = countFactor(A)
    res = 0
    for v in A:
        res += counter[v] == 1
    print(res)
