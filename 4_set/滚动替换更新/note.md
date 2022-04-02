```Python
class Solution:
    def solve(self, coins, quantities):
        """Return the number of distinct coin sum values you can get by using non-empty group of these coins."""
        res = set([0])
        for index, coin in enumerate(coins):
            cur = set()
            for pre in res:
                for count in range(quantities[index] + 1):
                    cur.add(pre + coin * count)
            res = cur

        return len(res) - 1

    def solve2(self, coins, quantities):
        dp = 1
        for coin, count in zip(coins, quantities):
            for _ in range(count):
                # 相当于集合并集操作
                dp |= dp << coin

        return bin(dp).count('1') - 1
```

本质上是 dp

```Python
class Solution:
    def solve(self, nums):
        res = set()
        cur = set()
        for num in nums:
            cur = {num | y for y in cur} | {num}
            res |= cur
        return len(res)

```
