# 某人又命名了一种取模数，这种数字必须满足各位数字之和 mod N 为 0。
# 现在大家又要玩游戏了，指定一个整数闭区间 [a.b]，问这个区间内有多少个取模数。


# 1≤a,b≤231−1,
# 1≤N<100


from functools import lru_cache


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, curSum: int, isLimit: bool) -> int:
        """当前在第pos位，和为curSum，isLimit表示是否贴合上界"""
        if pos == 0:
            return curSum % MOD == 0

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            res += dfs(pos - 1, curSum + cur, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True)


while True:
    try:
        left, right, MOD = map(int, input().split())
        print(cal(right) - cal(left - 1))
    except EOFError:
        break

