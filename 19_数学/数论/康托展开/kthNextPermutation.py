# 树状数组二分 O(nlogΣ)
# https://leetcode.cn/problems/smallest-palindromic-rearrangement-ii/

from collections import Counter, defaultdict, deque
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


def minAdjacentSwap(nums1: List[int], nums2: List[int]) -> int:
    """求使两个数组相等的最少邻位交换次数

    映射+求逆序对

    时间复杂度`O(nlogn)`
    """

    def countInversionPair(nums: List[int]) -> int:
        """计算逆序对的个数

        sortedList解法

        时间复杂度`O(nlogn)`
        """
        res = 0
        sl = SortedList()
        for num in reversed(nums):
            pos = sl.bisect_left(num)
            res += pos
            sl.add(num)
        return res

    # 含有重复元素的映射 例如nums [1,3,2,1,4] 表示已经排序的数组  [0,1,2,3,4]
    # 那么nums1 [1,1,3,4,2] 就 映射到 [0,3,1,4,2]
    mapping = defaultdict(deque)
    for index, num in enumerate(nums2):
        mapping[num].append(index)

    for index, num in enumerate(nums1):
        mapped = mapping[num].popleft()
        nums1[index] = mapped

    res = countInversionPair(nums1)

    return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/minimum-adjacent-swaps-to-reach-the-kth-smallest-number/
    # 测试用例是按存在第 k 个最小妙数而生成的。
    class Solution:
        def getMinSwaps(self, num: str, k: int) -> int:
            nums = [int(x) for x in num]
            res = kthNextPermutation(nums, k, inplace=False)
            if res is None:
                raise ValueError("不存在第k个最小妙数")
            return minAdjacentSwap(nums, res)

        # 3518. 最小回文排列 II
        # https://leetcode.cn/problems/smallest-palindromic-rearrangement-ii/description/
        def smallestPalindrome(self, s: str, k: int) -> str:
            if len(s) == 1:
                return "" if k > 1 else s
            ords = [ord(x) for x in s]
            kthLeftOrds = kthNextPermutation(sorted(ords[: len(ords) // 2]), k - 1)
            if kthLeftOrds is None:
                return ""
            kthLeft = "".join(chr(x) for x in kthLeftOrds)
            res = kthLeft
            if len(ords) % 2 == 1:
                res += chr(ords[len(ords) // 2])
            return res + kthLeft[::-1]

    assert Solution().getMinSwaps("5489355142", 4) == 2
