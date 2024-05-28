def min2(a: int, b: int) -> int:
    return a if a < b else b


class JosephusCircle:
    """O(min(n,k ln(n/k)))"""

    @staticmethod
    def survivor(n: int, k: int) -> int:
        """
        有n个人围成一圈,编号为0,1,2,...,n-1。
        从0开始报数,报到第k(k>=1)个人就出局, 然后从下一个人开始重新报数。
        求最后一个幸存者的编号。
        """
        if k == 1:
            return n - 1
        if k >= n:
            return JosephusCircle._survivorBF(n, k)
        survivor = JosephusCircle._survivorBF(k, k)
        i = k + 1
        while i <= n:
            t = (i - survivor - 1 + (k - 2)) // (k - 1)
            t = min2(t, n + 1 - i)
            i += t - 1
            survivor = (survivor + k * t) % i
            i += 1
        return survivor

    @staticmethod
    def dieAtRound(n: int, k: int, round: int) -> int:
        """
        求第round轮出局的人的编号。
        """
        alive = n - round + 1
        if k == 1:
            return round - 1
        who = (k - 1) % alive
        i = alive + 1
        while i <= n:
            t = (i - who - 1 + (k - 2)) // (k - 1)
            t = min2(t, n + 1 - i)
            i += t - 1
            who = (who + k * t) % i
            i += 1
        return who

    @staticmethod
    def dieTime(n: int, k: int, who: int) -> int:
        """
        求编号为who的人第几轮出局。
        O(min(n,kln(n/k))
        """
        if (who + 1) % k == 0:
            return (who + 1) // k
        turn = n // k
        if turn <= 1:
            return JosephusCircle._dieTimeBF(n, k, who)
        next = (who - turn * k) if (who >= turn * k) else (n + who - (who + 1) // k - turn * k)
        return JosephusCircle.dieTime(n - turn, k, next % n) + turn

    @staticmethod
    def _survivorBF(n: int, k: int) -> int:
        if n == 1:
            return 0
        res = 0
        for i in range(2, n + 1):
            res = (res + k) % i
        return res

    @staticmethod
    def _dieTimeBF(n: int, k: int, who: int) -> int:
        if (k - 1) % n == who:
            return 1
        return JosephusCircle._dieTimeBF(n - 1, k, (who - k) % n) + 1


# josephus
def josephus(n: int, jump: int, k: int) -> int:
    """
    约瑟夫环问题(从1开始)
    n: 总人数
    jump: 每次跳过的人数
    k: 第K个被点到的人(从1开始)
    O(jump*logn)
    """
    k = n - k
    if jump <= 1:
        return n - k
    i = k
    while i < n:
        r = (i - k + jump - 2) // (jump - 1)
        if (i + r) > n:
            r = n - i
        elif r == 0:
            r = 1
        i += r
        k = (k + (r * jump)) % i
    return k + 1


if __name__ == "__main__":
    # https://leetcode.cn/problems/yuan-quan-zhong-zui-hou-sheng-xia-de-shu-zi-lcof/solutions/607638/jian-zhi-offer-62-yuan-quan-zhong-zui-ho-dcow/
    class Solution:
        def iceBreakingGame(self, num: int, target: int) -> int:
            return JosephusCircle.survivor(num, target)
