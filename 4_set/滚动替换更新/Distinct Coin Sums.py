class Solution:
    def solve(self, coins, quantities):
        """Return the number of distinct coin sum values you can get by using non-empty group of these coins."""
        dp = set([0])
        for index, coin in enumerate(coins):
            ndp = set()
            for pre in dp:
                for count in range(quantities[index] + 1):
                    ndp.add(pre + coin * count)
            dp = ndp

        return len(dp) - 1

    def solve2(self, coins, quantities):
        dp = 1
        for coin, count in zip(coins, quantities):
            for _ in range(count):
                # 相当于集合并集操作
                dp |= dp << coin

        return bin(dp).count('1') - 1


print(Solution().solve([4, 2, 1], [1, 2, 1]))
# We can have the following distinct coin sums

# [1]
# [2]
# [1, 2]
# [4]
# [1, 4]
# [2, 4]
# [1, 2, 4]
# [2, 2, 4]
# [1, 2, 2, 4]
