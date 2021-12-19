from collections import Counter

# 对于每个数位 d ，这个数位 恰好 在 x 中出现 d 次。那么整数 x 就是一个 数值平衡数 。
# 请你返回 严格大于 n 的 最小数值平衡数

# 0 <= n <= 106
class Solution:
    def nextBeautifulNumber(self, n: int) -> int:
        def check(num: int) -> bool:
            counter = Counter(list(str(num)))
            return all([int(k) == v for k, v in counter.items()])

        res = n + 1
        while not check(res):
            res += 1
        return res

    def nextBeautifulNumber2(self, n: int) -> int:
        return next(
            n for n in range(n + 1, 1234567) if all(int(x) == t for x, t in Counter(str(n)).items())
        )


print(Solution().nextBeautifulNumber(n=1))
# 输出：22
# 解释：
# 22 是一个数值平衡数，因为：
# - 数字 2 出现 2 次
# 这也是严格大于 1 的最小数值平衡数。
