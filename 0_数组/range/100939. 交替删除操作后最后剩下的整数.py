# 100939. 交替删除操作后最后剩下的整数
# https://leetcode.cn/problems/last-remaining-integer-after-alternating-deletion-operations/
# 给你一个整数 n。
# 我们将 1 到 n 的整数按从左到右的顺序排成一个序列。然后，交替 地执行以下两种操作，直到只剩下一个整数为止，从操作 1 开始：
# 操作 1：从左侧开始，隔一个数删除一个数。
# 操作 2：从右侧开始，隔一个数删除一个数。
# 返回最后剩下的那个整数。


class Solution:
    def lastInteger(self, n: int) -> int:
        seq = range(1, n + 1)
        while len(seq) > 1:
            seq = seq[::2][::-1]
        return seq.start
