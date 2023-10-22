# 100097. 合法分组的最少组数
# https://leetcode.cn/problems/minimum-number-of-groups-to-create-a-valid-assignment/description/
#
# 给你一个长度为 n 下标从 0 开始的整数数组 nums 。
#
# 我们想将下标进行分组，使得 [0, n - 1] 内所有下标 i 都 恰好 被分到其中一组。
#
# 如果以下条件成立，我们说这个分组方案是合法的：
#
# 对于每个组 g ，同一组内所有下标在 nums 中对应的数值都相等。
# 对于任意两个组 g1 和 g2 ，两个组中 下标数量 的 差值不超过 1 。
# 请你返回一个整数，表示得到一个合法分组方案的 最少 组数。
#
# 最后每种频率需要拆成size和size+1两种
# !频率的种类数不超过根号n，因此可以直接枚举size

from typing import List, Tuple
from collections import Counter


class Solution:
    def minGroupsForValidAssignment(self, nums: List[int]) -> int:
        freqCounter = Counter(Counter(nums).values())
        res = len(nums)
        for size in range(1, len(nums) + 1):
            tmp = 0
            for num in freqCounter:
                count1, count2, ok = splitToKAndKPlusOne(num, size)
                if not ok:
                    break
                tmp += (count1 + count2) * freqCounter[num]
            else:
                res = min(res, tmp)
        return res


def splitToKAndKPlusOne(num: int, k: int, minimize=True) -> Tuple[int, int, bool]:
    """将 num 拆分成 k 和 k+1 的和，使得拆分的个数最(多/少).

    Args:
        num (int): 正整数.
        k (int): 正整数.
        minimize (bool, optional): 是否使得拆分的个数最少. 默认为最少(true).

    Returns:
        Tuple[int, int, bool]: count1和count2分别是拆分成k和k+1的个数,ok表示是否可以拆分.
    """
    if minimize:
        count2 = (num + k) // (k + 1)
        diff = (k + 1) * count2 - num
        if diff > count2:
            return 0, 0, False
        return diff, count2 - diff, True
    count1 = num // k
    diff = num - k * count1
    if diff > count1:
        return 0, 0, False
    return count1 - diff, diff, True
