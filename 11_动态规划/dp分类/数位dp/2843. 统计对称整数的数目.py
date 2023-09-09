# 2843. 统计对称整数的数目
# https://leetcode.cn/problems/count-symmetric-integers/

# 给你两个正整数 low 和 high 。
# 对于一个由 2 * n 位数字组成的整数 x ，
# !如果其前 n 位数字之和与后 n 位数字之和相等，则认为这个数字是一个对称整数。
# 返回在 [low, high] 范围内的 对称整数的数目 。

from functools import lru_cache

MOD = int(1e9 + 7)


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(
        pos: int,
        hasLeadingZero: bool,
        isLimit: bool,
        halfLen: int,
        diff: int,
    ) -> int:
        """
        当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界.
        halfLen表示前一半的长度,diff表示前一半的和与后一半的和的差值
        """
        if pos == n:
            return int(halfLen > 0 and diff == 0)

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, True, (isLimit and cur == up), halfLen, diff)
            else:
                nextHalfLen = halfLen
                if halfLen == 0:
                    allLen = n - pos
                    if allLen & 1:
                        continue
                    nextHalfLen = allLen // 2
                isInPreHalf = pos < n - nextHalfLen
                nextDiff = diff + (cur if isInPreHalf else -cur)
                res += dfs(pos + 1, False, (isLimit and cur == up), nextHalfLen, nextDiff)
        return res % MOD

    nums = list(map(int, str(upper)))
    n = len(nums)
    return dfs(0, True, True, 0, 0) % MOD


class Solution:
    def countSymmetricIntegers(self, low: int, high: int) -> int:
        return (cal(high) - cal(low - 1)) % MOD


if __name__ == "__main__":
    print(Solution().countSymmetricIntegers(1, 100))
