# https://maspypy.github.io/library/other/solve_hukumenzan.hpp
# 覆面算(口算难题/竖式计算)
# a + b = c を解く

from typing import List, Tuple


def solveHukumenzan(A: str, B: str, C: str) -> List[Tuple[int, int, int]]:
    res = []
    v = list(set(list(A + B + C)))
    if len(v) > 10:
        return res
    order = list(range(10))
    while True:
        ok, order = nextPermutation(order, inPlace=True)
        if not ok:
            break
        mp = {v[i]: order[i] for i in range(len(v))}
        if mp[A[0]] == 0 or mp[B[0]] == 0 or mp[C[0]] == 0:
            continue
        a, b, c = 0, 0, 0
        for x in A:
            a = 10 * a + mp[x]
        for x in B:
            b = 10 * b + mp[x]
        for x in C:
            c = 10 * c + mp[x]
        if a + b == c:
            res.append((a, b, c))
    return res


from typing import Any, MutableSequence, Tuple


def nextPermutation(nums: MutableSequence[Any], inPlace=False) -> Tuple[bool, MutableSequence[Any]]:
    """返回下一个字典序的排列，如果不存在，返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] >= nums[left]:  # 1. 找到最后一个递增位置
        left -= 1
    if left == 0:  # 全部递减
        return False, nums
    last = left - 1  # 最后一个递增位置

    while nums[right] <= nums[last]:  # 2. 找到最小的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递减数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


# https://atcoder.jp/contests/abc198/tasks/abc198_d
if __name__ == "__main__":
    # A,B,C都是小写字母组成的长度不超过10的字符串

    A, B, C = input(), input(), input()
    res = solveHukumenzan(A, B, C)
    if not res:
        print("UNSOLVABLE")
    else:
        a, b, c = res[0]
        print(a)
        print(b)
        print(c)
