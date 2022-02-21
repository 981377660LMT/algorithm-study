# 一堆个数为n的石子，Alice和Bob轮流取。
# Alice一次能取[1,p]个石子，牛妹一次能取[1,q]个石子。
# 拿到最后一个石子的人赢。
class Solution:
    def Gameresults(self, n, p, q):
        # write code here
        if n <= p:
            return 1
        if p == q:
            mod = n % (p + 1)
            if mod == 0:
                return -1
            else:
                return 1
        if p > q:
            return 1
        return -1
