# 最大化车的得分
# 在棋盘上放两个车,让他们不互相攻击 (攻击是指两个车在同一行或同一列上)。
# 求最大的分
class Solution:
    def solve(self, board):
        matrix = [(val, (r, c)) for r, row in enumerate(board) for c, val in enumerate(row)]
        max1 = max(matrix)
        max2 = max(v for v in matrix if v != max1)  # 直接选出来不等的最大的
        cands = [max1, max2]  # 两个就可以取所有情况,一个不行

        res = 0
        for v, (r, c) in cands:
            for nv, (nr, nc) in matrix:
                if nr != r and c != nc:
                    res = max(res, v + nv)
        return res
