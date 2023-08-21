# 如果一个十进制数能够被它的各位数字之和整除，则称这个数为“月之数”。
# 给定整数 L 和 R，你需要计算闭区间 [L,R] 中有多少个“月之数”。
# 1≤L,R<2^31

from functools import lru_cache


# !”求出数字的和，去判断各位数字之和能否整除他”的想法是错误的，因为每个数字都要判断一次，记搜变成暴搜。
# 转换思路，能够整除就是模数为0，所以可以考虑把每步的模数给算出来，算到最后一位判断模数是否为0即可


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, isLimit: bool, allSum: int, mod: int, curSum: int) -> int:
        """当前在第pos位，isLimit表示是否贴合上界,sum表示各位数字之和,mod表示原数字模各位数字之和的模数"""
        if pos == 0:
            return mod == 0 and curSum == allSum

        res = 0
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            res += dfs(
                pos - 1, (isLimit and cur == up), allSum, (mod * 10 + cur) % allSum, curSum + cur
            )
        return res

    nums = []
    while upper:
        nums.append(upper % 10)
        upper //= 10
    # 枚举各位数字之和
    return sum(dfs(len(nums), True, sum, 0, 0) for sum in range(1, 9 * len(nums)))


left, right = map(int, input().split())
print(cal(right) - cal(left - 1))
