# https://oi-wiki.org/misc/job-order/

# 在一台机器上规划任务(排序不等式，接水问题)
# 你有 n 个任务，要求你找到一个代价最小的的顺序执行他们。
# 第 i 个任务花费的时间是 t_i，而第 i 个任务等待 t 的时间会花费 f_i(t) 的代价。

# 0.相同的单增函数(排队打水问题，排序不等式)
# !按照ti从小到大排序
# 1.线性代价函数，考虑微扰后的变换情况，贪心地选取最优解
# !fi(x) = cix+di => 按照 ci/ti 从小到大排序
# 2. 指数代价函数，考虑微扰后的变换情况，贪心地选取最优解
# !fi(x) = ci*e^(a*x) => 按照 (1-e^(a*ti))/ci 从小到大排序


# !线性代价函数


if __name__ == "__main__":
    # 100391. 对 Bob 造成的最少伤害
    # https://leetcode.cn/contest/biweekly-contest-138/problems/minimum-amount-of-damage-dealt-to-bob/
    # Bob 有 n 个敌人，如果第 i 个敌人还活着（也就是健康值 health[i] > 0 的时候），每秒钟会对 Bob 造成 damage[i] 点 伤害。
    # 每一秒中，在敌人对 Bob 造成伤害 之后 ，Bob 会选择 一个 还活着的敌人进行攻击，该敌人的健康值减少 power 。
    # 请你返回 Bob 将 所有 n 个敌人都消灭之前，最少 会收到多少伤害。

    from typing import List
    from functools import cmp_to_key

    class Solution:
        def minDamage(self, power: int, damage: List[int], health: List[int]) -> int:
            def cmp(i0: int, i1: int) -> int:
                return -(damage[i0] * times[i1] - damage[i1] * times[i0])

            times = [(v + power - 1) // power for v in health]
            order = sorted(range(len(damage)), key=cmp_to_key(cmp))  # !排队顺序

            res, curSum = 0, sum(damage)
            for id in order:
                res += curSum * times[id]
                curSum -= damage[id]
            return res
