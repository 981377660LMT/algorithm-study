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
            if cand == target:
                permCount = target
                i += 1
                pre = sl[i]
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
    # https://leetcode.cn/problems/minimum-adjacent-swaps-to-reach-the-kth-smallest-number/
    # num = "11112", k = 4
    nums = [1, 1, 1, 1, 2]
    print(kthNextPermutation(nums, 2, prev=False, inplace=False))
