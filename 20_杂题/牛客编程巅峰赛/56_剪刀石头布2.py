# Alice和Bob打牌，每人都有n张牌
# Alice的牌里有p1张石头牌，q1张剪刀牌，m1张布牌。
# Bob的牌里有p2张石头牌，q2张剪刀牌，m2张布牌。
# Alice 获胜得一分，失败扣一分，平局不得分也不扣分
# Alice知道Bob每次要出什么牌，请你安排策略，使Alice的分最多。

#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
#
# @param n int
# @param p1 int
# @param q1 int
# @param m1 int
# @param p2 int
# @param q2 int
# @param m2 int
# @return int
#
class Solution:
    def Highestscore(self, n, p1, q1, m1, p2, q2, m2):
        # write code here
        # 最高得分=牛牛最多赢的局-牛妹最少赢的局
        c1 = min(p1, q2)
        p1 -= c1
        q2 -= c1
        c2 = min(q1, m2)
        q1 -= c2
        m2 -= c2
        c3 = min(m1, p2)
        m1 -= c3
        p2 -= c3

        # 最多赢得据
        aliceWin = c1 + c2 + c3
        res = aliceWin

        tie = min(p1, p2) + min(q1, q2) + min(m1, m2)
        res -= n - (aliceWin + tie)
        return res
