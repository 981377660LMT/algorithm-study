# 下一个排列/上一个排列


from typing import Any, List, Tuple


def nextPermutation(nums: List[Any], inPlace=False) -> Tuple[bool, List[Any]]:
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


def prePermutation(nums: List[Any], inPlace=False) -> Tuple[bool, List[Any]]:
    """返回前一个字典序的排列,如果不存在,返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] <= nums[left]:  # 1. 找到最后一个递减位置
        left -= 1
    if left == 0:  # 全部递增
        return False, nums
    last = left - 1  # 最后一个递减位置

    while nums[right] >= nums[last]:  # 2. 找到最大的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递增数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


from sortedcontainers import SortedSet


# https://maspypy.github.io/library/seq/kth_next_permutation.hpp
def kthNextPermutationDistinct(
    unique: List[Any], k: int, inPlace=False, prev=False
) -> Tuple[bool, List[Any], int]:
    """下k个字典序的排列(无重复元素)

    Args:
        unique (List[Any]): `无重复元素`的数组
        k (int): 后续第k个(`本身算第0个`)
        inPlace (bool, optional): 是否原地修改. 默认为False
        prev (bool, optional): 使用next还是prev. 默认使用next

    Returns:
        Tuple[bool, List[Any], int]: `是否存在, 下k个排列, 需要移动的元素个数`
    """
    if not inPlace:
        unique = unique[:]
    rank, q = [], []
    ss = SortedSet() if not prev else SortedSet(key=lambda x: -x)
    while k and unique:
        n = len(rank) + 1
        p = unique[-1]
        now = ss.bisect_left(p)
        k += now
        r = k % n
        k //= n
        rank.append(r)
        q.append(unique[-1])
        ss.add(unique[-1])
        unique.pop()

    if k:
        return False, unique, len(rank)

    move = len(rank)
    while len(rank):
        r = rank.pop()
        it = ss[r]
        unique.append(it)
        ss.remove(it)
    return True, unique, move


from collections import Counter
from typing import List, Optional
from sortedcontainers import SortedList


def kthNextPermutation(nums: List[int], k: int, inplace=False, prev=False) -> Optional[List[int]]:
    """下k个字典序的排列(可以存在重复元素)

    Args:
        nums: 原排列数组.
        k (int): 后续第k个(`本身算第0个`)字典序的排列.
        inplace (bool, optional): 是否原地修改. 默认为False.
        prev (bool, optional): 使用next还是prev. 默认使用next.

    Returns:
        Optional[List[int]]: 下k个字典序的排列,如果不存在,返回None.
    """
    if not nums:
        return
    counter = Counter([nums[-1]])
    sl = SortedList([nums[-1]]) if not prev else SortedList([nums[-1]], key=lambda x: -x)

    fac = 1
    facPtr = 1
    curPerm = 0
    overlap = 1  # 重复元素的个数的乘积
    allPerm = 1  # 后缀里的所有排列个数
    for right in range(len(nums) - 2, -1, -1):
        if curPerm + k < allPerm:
            break
        num = nums[right]
        counter[num] += 1
        overlap *= counter[num]

        smaller = 0
        pos = sl.bisect_left(num)
        if pos == len(sl) or sl[pos] != num:  # set去重
            sl.add(num)
        for pre in sl.islice(0, pos):
            smaller += (fac * counter[pre]) // overlap

        facPtr += 1
        fac *= facPtr
        curPerm += smaller
        allPerm = fac // overlap

    if curPerm + k >= allPerm:
        return

    res = []
    fac //= facPtr
    permCount = 0
    target = curPerm + k
    while permCount != target:
        for i, pre in enumerate(sl):
            curPerm = (fac * counter[pre]) // overlap  # 以当前元素开头的排列个数
            cand = permCount + curPerm
            if cand <= target:
                permCount = cand
                continue
            facPtr -= 1
            fac //= facPtr
            overlap //= counter[pre]
            res.append(pre)
            counter[pre] -= 1
            if not counter[pre]:
                sl.pop(i)
            break

    for pre in sl:
        res += [pre] * counter[pre]
    if inplace:
        nums[-len(res) :] = res
        return nums
    return nums[: len(nums) - len(res)] + res


if __name__ == "__main__":
    isOk, nextP = nextPermutation([1, 2, 3])
    if isOk:
        print("nextP:", nextP)
    isOk, nextP = nextPermutation(list(map(int, "12345")))
    if isOk:
        print("nextP:", nextP)

    isOk, preP = prePermutation([3, 2, 1])
    if isOk:
        print("preP:", preP)
    isOk, preP = prePermutation(list(map(int, "12345")))
    if isOk:
        print("preP:", preP)

    isOk, kthP, move = kthNextPermutationDistinct(list(map(int, "12345")), 5)
    if isOk:
        print("kthP:", kthP, "move:", move)

    a1 = [1, 2, 3, 4, 5, 6, 7, 9, 23, 14, 56, 99, 876, 222, 444, 555]
    a2 = [1, 2, 3, 4, 5, 6, 7, 9, 23, 14, 56, 99, 876, 222, 444, 555]

    k = 999
    for _ in range(k):
        isOk, a1 = nextPermutation(a1, inPlace=True)
        if not isOk:
            break
    assert a1 == kthNextPermutationDistinct(a2, k)[1] == kthNextPermutation(a2, k)

    nums1 = [5, 4, 3, 2, 1, 10, 9, 8, 7, 6, 11, 12, 13, 14, 15, 16, 16]
    nums2 = nums1[:]
    k = 200
    for i in range(k):
        isOk, nums1 = prePermutation(nums1)
        if not isOk:
            break
        assert nums1 == kthNextPermutation(nums2, i + 1, prev=True)
