# Alice和Bob打牌，每人都有n张牌
# Alice的牌里有p1张石头牌，q1张剪刀牌，m1张布牌。
# Bob的牌里有p2张石头牌，q2张剪刀牌，m2张布牌。
# Alice知道Bob每次要出什么牌，请你安排策略，使Alice获胜次数最多。
# 输出获胜次数。
#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# ​请返回Alice能赢的最多局数
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
    def Mostvictories(self, n, p1, q1, m1, p2, q2, m2):
        # write code here
        return min(p1, q2) + min(q1, m2) + min(m1, p2)

