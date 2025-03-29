# 3494. 酿造药水需要的最少总时间-流水线作业
# https://leetcode.cn/problems/find-the-minimum-amount-of-time-to-brew-potions/description/
# 给你两个长度分别为 n 和 m 的整数数组 skill 和 mana 。
# 在一个实验室里，有 n 个巫师，他们必须按顺序酿造 m 个药水。每个药水的法力值为 mana[j]，并且每个药水 必须 依次通过 所有 巫师处理，才能完成酿造。第 i 个巫师在第 j 个药水上处理需要的时间为 timeij = skill[i] * mana[j]
# 由于酿造过程非常精细，药水在当前巫师完成工作后 必须 立即传递给下一个巫师并开始处理。这意味着时间必须保持 同步，确保每个巫师在药水到达时 马上 开始工作。
# 返回酿造所有药水所需的 最短 总时间。
#
# !凸壳二分
# https://leetcode.cn/problems/find-the-minimum-amount-of-time-to-brew-potions/solutions/3625002/onm-log-ntu-ke-er-fen-by-hqztrue-ktyh/?slug=find-the-minimum-amount-of-time-to-brew-potions&region=local_v2

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def minTime(self, skill: List[int], mana: List[int]) -> int:
        n = len(skill)
        lastFinish = [0] * n
        for m in mana:
            sum_ = 0
            for x, last in zip(skill, lastFinish):
                sum_ = max2(sum_, last) + x * m
            lastFinish[-1] = sum_
            for i in range(n - 2, -1, -1):
                lastFinish[i] = lastFinish[i + 1] - m * skill[i + 1]
        return lastFinish[-1]
