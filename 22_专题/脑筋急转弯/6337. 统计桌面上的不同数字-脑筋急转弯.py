# 给你一个正整数 n ，开始时，它放在桌面上。
# 在 1e9 天内，每天都要执行下述步骤：

# 对于出现在桌面上的每个数字 x ，找出符合 1 <= i <= n 且满足 x % i == 1 的所有数字 i 。
# 然后，将这些数字放在桌面上。
# 返回在 1e9 天之后，出现在桌面上的 不同 整数的数目。


# !脑筋急转弯,注意特判n=1的情况
# 由于 n%(n-1)==1 因此第一天后n-1会出现在桌面上
# 第二天后n-2会出现在桌面上
# ...
# 这样最多n天就能把2到n放在桌面上
# 由于1永远不会出现在桌面上,因此最多n-1个不同的整数出现在桌面上


class Solution:
    def distinctIntegers(self, n: int) -> int:
        if n == 1:
            return 1
        return n - 1

    def distinctIntegers2(self, n: int) -> int:
        """模拟,直到不发生变化结束"""
        dp = set([n])
        while True:
            ndp = dp.copy()
            for x in dp:
                for i in range(1, n + 1):
                    if x % i == 1:
                        ndp.add(i)
            if len(ndp) == len(dp):
                break
            dp = ndp
        return len(dp)
