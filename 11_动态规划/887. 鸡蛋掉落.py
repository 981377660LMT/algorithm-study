# 887. 鸡蛋掉落
# https://leetcode.cn/problems/super-egg-drop/
# 给你 k 枚相同的鸡蛋，并可以使用一栋从第 1 层到第 n 层共有 n 层楼的建筑。
# 已知存在楼层 f ，满足 0 <= f <= n ，任何从 高于 f 的楼层落下的鸡蛋都会碎，从 f 楼层或比它低的楼层落下的鸡蛋都不会破。
# 每次操作，你可以取一枚没有碎的鸡蛋并把它从任一楼层 x 扔下（满足 1 <= x <= n）。
# 如果鸡蛋碎了，你就不能再次使用它。如果某枚鸡蛋扔下后没有摔碎，则可以在之后的操作中 重复使用 这枚鸡蛋。
# 请你计算并返回要确定 f 确切的值 的 最小操作次数 是多少？
#
# !一次扔蛋可以“分裂”成两个子问题，再加上当前检测的那一层。
#
# O(logn)解决任意鸡蛋个数的问题
# https://leetcode.cn/problems/egg-drop-with-2-eggs-and-n-floors/solutions/2948790/olognjie-jue-ren-yi-ji-dan-ge-shu-de-wen-yjlq/


from functools import lru_cache


@lru_cache(None)
def dfs(i: int, j: int):
    """有i次操作机会和j个鸡蛋的情况下,最多可以检测多少层楼."""
    if i == 0 or j == 0:
        return 0
    return dfs(i - 1, j - 1) + dfs(i - 1, j) + 1


class Solution:
    def superEggDrop(self, k: int, n: int) -> int:
        res = 1
        while dfs(res, k) < n:
            res += 1
        return res
