# 某人命名了一种不降数，这种数字必须满足从左到右各位数字呈非下降关系，如 123，446。
# 现在大家决定玩一个游戏，指定一个整数闭区间 [a,b]，问这个区间内有多少个不降数。
# # 1≤X≤Y≤231−1,


from functools import lru_cache


@lru_cache(None)
def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, pre: int, isLimit: bool) -> int:
        """当前在第pos位，前一个数为pre，isLimit表示是否贴合上界"""
        if pos == 0:
            return 1

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if cur < pre:
                continue
            res += dfs(pos - 1, cur, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True)


while True:
    try:
        left, right = map(int, input().split())
        print(cal(right) - cal(left - 1))
    except EOFError:
        break

