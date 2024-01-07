# 100163. 统计强大整数的数目
# https://leetcode.cn/problems/count-the-number-of-powerful-integers/submissions/493818546/
# 给你三个整数 start ，finish 和 limit 。同时给你一个下标从 0 开始的字符串 s ，表示一个 正 整数。
# 如果一个 正 整数 x 末尾部分是 s （换句话说，s 是 x 的 后缀），
# 且 x 中的每个数位至多是 limit ，那么我们称 x 是 强大的 。
# 请你返回区间 [start..finish] 内强大整数的 总数目 。
# 如果一个字符串 x 是 y 中某个下标开始（包括 0 ），到下标为 y.length - 1 结束的子字符串，
# 那么我们称 x 是 y 的一个后缀。比方说，25 是 5125 的一个后缀，但不是 512 的后缀。


from functools import lru_cache


def min2(a, b):
    return a if a < b else b


class Solution:
    def numberOfPowerfulInt(self, start: int, finish: int, limit: int, s: str) -> int:
        def cal(upper: int) -> int:
            @lru_cache(None)
            def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, matchPos: int) -> int:
                """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界,matchPos表示当前匹配到s的第几位."""
                if pos >= n:
                    return 1 if not hasLeadingZero else 0

                res = 0
                up = min2(nums[pos], limit) if isLimit else limit
                shouldMatch = pos >= offset
                for cur in range(up + 1):
                    if not shouldMatch:
                        if hasLeadingZero and cur == 0:
                            res += dfs(pos + 1, True, (isLimit and cur == nums[pos]), matchPos)
                        else:
                            res += dfs(
                                pos + 1,
                                False,
                                (isLimit and cur == nums[pos]),
                                matchPos,
                            )
                    else:
                        if cur == sNums[matchPos]:
                            res += dfs(
                                pos + 1,
                                False,
                                (isLimit and cur == nums[pos]),
                                matchPos + 1,
                            )
                return res

            if upper < int(s):
                return 0
            nums = list(map(int, str(upper)))
            sNums = list(map(int, s))
            n = len(nums)
            offset = max(0, len(nums) - len(s))
            return dfs(0, True, True, 0)

        return cal(finish) - cal(start - 1)
