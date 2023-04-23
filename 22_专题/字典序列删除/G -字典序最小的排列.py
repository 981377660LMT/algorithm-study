# https://atcoder.jp/contests/abc299/tasks/abc299_g
# G -字典序最小的排列
# 从序列 M 个数中顺序选出 N 个不同的数, 使得这 N 个数的字典序最小。


from typing import List
from collections import Counter


def minimumPermutation(nums: List[int], m: int) -> List[int]:
    remain = Counter(nums)
    visited = set()
    stack = []
    for num in nums:
        remain[num] -= 1
        if num in visited:
            continue
        while stack and stack[-1] > num and remain[stack[-1]] > 0:
            visited.remove(stack.pop())
        stack.append(num)
        visited.add(num)
    return stack[:m]


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))
    res = minimumPermutation(nums, m)
    print(*res, sep=" ")
