# 翻转灯泡/翻转01串
# [开关问题]01串翻转全变为零

# 有一个01串，长度为len(1<=len<=1e5)，我们可以对这个字符串进行翻转操作，定义如下：
# 我们可以对字符串任意位置进行操作，操作后，
# !`此位置与其相邻的两个位置`的字符改变(‘1’变’0’，‘0’变’1’)。
# 问最少翻转多少次这个字符串可以全部变为为’0’。
# 输出最少步数，如果无法实现，输出-1

# !处理环(有后效性的问题):分类讨论,是否要翻转第一个开关
# i从第1位至倒数第二位遍历，若该位为1,则flip(i+1), 若最后剩下00...00001，则不行。

from typing import List


INF = int(4e18)


def flip(nums: List[int], i: int) -> None:
    for i in range(i - 1, i + 2):
        if 0 <= i < len(nums):
            nums[i] ^= 1


def solve(s: str) -> int:
    def cal1(nums: List[int]) -> int:
        """按下第一个开关"""
        flip(nums, 0)
        res = 1
        for i in range(n - 1):
            if nums[i] == 1:
                flip(nums, i + 1)
                res += 1
        return res if nums[-1] == 0 else INF

    def cal2(nums: List[int]) -> int:
        """不按下第一个开关"""
        res = 0
        for i in range(n - 1):
            if nums[i] == 1:
                flip(nums, i + 1)
                res += 1
        return res if nums[-1] == 0 else INF

    n = len(s)
    cand1 = cal1(list(map(int, s)))
    cand2 = cal2(list(map(int, s)))
    res = min(cand1, cand2)
    return res if res != INF else -1


if __name__ == "__main__":
    assert solve("101") == 2
    assert solve("000") == 0
    assert solve("01") == -1
    assert solve("0001") == 3
    assert solve("10000") == -1
