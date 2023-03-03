# 这里有 n 个一样的骰子，每个骰子上都有 k 个面，分别标号为 1 到 k 。
# 给定三个整数 n ,  k 和 target ，
# 返回可能的方式(从总共 k^n 种方式中)滚动骰子的数量，
# 使正面朝上的数字之和等于 target 。
# 答案可能很大，你需要对 1e9 + 7 取模 。


class Solution:
    def numRollsToTarget(self, n: int, k: int, target: int) -> int:
        ...
