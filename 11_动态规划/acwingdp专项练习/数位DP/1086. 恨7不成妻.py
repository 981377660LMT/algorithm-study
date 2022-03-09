# 如果一个整数符合下面三个条件之一，那么我们就说这个整数和 7 有关：

# 整数中某一位是 7；
# 整数的每一位加起来的和是 7 的整数倍；
# 这个整数是 7 的整数倍。

# 对于每组数据，请计算 [L,R] 中和 7 无关的数字的平方和，并将结果对 109+7 取模后输出。

# 1≤L≤R≤1018

MOD = int(1e9) + 7

from functools import lru_cache

path = []  # path不可以记忆化，所以放在外面

# 答案不对
def cal(upper: int) -> int:
    if upper == 0:
        return 0

    @lru_cache(None)
    def dfs(pos: int, curSum: int, curNum: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """
        当前在第pos位，每一位加起来和为curSum，整数模7的值为mod7,
        hasLeadingZero表示前面都选的是0(即有没有开始选)
        isLimit表示是否贴合上界
        """
        if pos == 0:
            if not path:
                return 0
            return (int(''.join(path)) ** 2) % MOD

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            if cur == 7:
                continue
            # 还没开始选
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, curSum, curNum, True, isLimit)
            # 开始选了
            else:
                path.append(str(cur))
                res += dfs(
                    pos - 1,
                    (curSum + cur) % 7,
                    (curNum * 10 + cur) % 7,
                    False,
                    (isLimit and cur == up),
                )
                res %= MOD
                path.pop()
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, 0, True, True)


n = int(input())
for _ in range(n):
    left, right = map(int, input().split())
    print(cal(right) - cal(left - 1))
