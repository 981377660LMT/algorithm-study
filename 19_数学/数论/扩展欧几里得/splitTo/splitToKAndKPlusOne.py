from collections import Counter
from typing import List, Tuple


def splitToKAndKPlusOne(num: int, k: int, minimize=True) -> Tuple[int, int, bool]:
    """将 num 拆分成 k 和 k+1 的和，使得拆分的个数最(多/少).

    Args:
        num (int): 正整数.
        k (int): 正整数.
        minimize (bool, optional): 是否使得拆分的个数最少. 默认为最少(true).

    Returns:
        Tuple[int, int, bool]: count1和count2分别是拆分成k和k+1的个数，ok表示是否可以拆分.
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


if __name__ == "__main__":
    # 2870. 使数组为空的最少操作次数
    # https://leetcode.cn/problems/minimum-number-of-operations-to-make-array-empty/
    class Solution:
        def minOperations(self, nums: List[int]) -> int:
            counter = Counter(nums)
            res = 0
            for v in counter.values():
                c1, c2, ok = splitToKAndKPlusOne(v, 2)
                if not ok:
                    return -1
                res += c1 + c2
            return res

    def checkWithBruteForce(num: int, k: int, minimize=True) -> Tuple[int, int, bool]:
        res = [0, 0, False]
        for count1 in range(num + 1):
            for count2 in range(num + 1):
                sum_ = count1 * k + (count2) * (k + 1)
                if sum_ == num:
                    if not res[2] or minimize ^ ((count1 + count2) > (res[0] + res[1])):
                        res = [count1, count2, True]
        return res

    for num in range(1, 50):
        for k in range(1, 50):
            for minimize in [True, False]:
                assert splitToKAndKPlusOne(num, k, minimize) == checkWithBruteForce(
                    num, k, minimize
                )

    print("pass")
