# Windy 定义了一种 Windy 数：不含前导零且相邻两个数字之差至少为 2 的正整数被称为 Windy 数。
# Windy 想知道，在 A 和 B 之间，包括 A 和 B，总共有多少个 Windy 数？

# 1≤A≤B≤2×109


from functools import lru_cache


def cal(upper: int) -> int:
    if upper == 0:
        return 1

    @lru_cache(None)
    def dfs(pos: int, pre: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位，前一个数为pre，isLimit表示是否贴合上界"""
        """此题需要记录前导零hasLeadingZero"""
        if pos == 0:
            return 1

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if abs(pre - cur) < 2:
                continue
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, -2, True, (isLimit and cur == up))
            else:
                res += dfs(pos - 1, cur, False, (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), -10, True, True)


# print(cal(1))
print(cal(10))
# while True:
#     try:
#         left, right = map(int, input().split())
#         print(cal(right) - cal(left - 1))
#     except EOFError:
#         break

