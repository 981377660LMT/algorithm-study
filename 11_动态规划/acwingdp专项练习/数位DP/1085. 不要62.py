# 杭州人称那些傻乎乎粘嗒嗒的人为 62（音：laoer）。
# 不吉利的数字为所有含有 4 或 62 的号码
# 你的任务是，对于每次给出的一个牌照号区间 [n,m]，求出有多少个不含4和62的数。

# 1≤n≤m≤109
# 1≤N<100


from functools import lru_cache


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, pre: int, isLimit: bool) -> int:
        """当前在第pos位，前一个数pre，isLimit表示是否贴合上界"""
        if pos == 0:
            return 1

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if cur == 4:
                continue
            if pre == 6 and cur == 2:
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
    left, right = map(int, input().split())
    if left == right == 0:
        break
    print(cal(right) - cal(left - 1))

