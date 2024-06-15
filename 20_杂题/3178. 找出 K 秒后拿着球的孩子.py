# 递枕头/传球问题
# 3178. 找出 K 秒后拿着球的孩子
# https://leetcode.cn/problems/find-the-child-who-has-the-ball-after-k-seconds/description/
# 给你两个 正整数 n 和 k。有 n 个编号从 0 到 n - 1 的孩子按顺序从左到右站成一队。
# 最初，编号为 0 的孩子拿着一个球，并且向右传球。每过一秒，拿着球的孩子就会将球传给他旁边的孩子。一旦球到达队列的 任一端 ，即编号为 0 的孩子或编号为 n - 1 的孩子处，传球方向就会 反转 。
# 返回 k 秒后接到球的孩子的编号。


class Solution:
    def numberOfChild(self, n: int, k: int) -> int:
        div, mod = divmod(k, n - 1)
        return n - mod - 1 if div & 1 else mod
