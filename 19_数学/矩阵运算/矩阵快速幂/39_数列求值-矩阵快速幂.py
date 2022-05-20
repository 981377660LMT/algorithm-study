MOD = int(1e9 + 7)
Matrix = list[list[int]]
#
# 输出序列的第n项
# @param n long长整型 序列的项数
# @param b long长整型 系数
# @param c long长整型 系数
# @return long长整型
#
# https://www.nowcoder.com/practice/3dfb8d459a2e439bb041c2503d14e5c2?tpId=185&tqId=35110&rp=1&ru=/ta/weeklycontest-history&qru=/ta/weeklycontest-history&difficulty=&judgeStatus=&tags=/question-ranking
class Solution:
    def nthElement(self, n: int, b: int, c: int) -> int:
        """求an模int(1e9+7)"""

        def multi(a: Matrix, b: Matrix) -> Matrix:
            """matrix multi"""
            res = [[0, 0], [0, 0]]
            for r in range(2):
                for c in range(2):
                    for k in range(2):
                        res[r][c] += a[r][k] * b[k][c]
                        res[r][c] %= MOD
            return res

        def qpow(matrix: Matrix, k: int) -> Matrix:
            res = [[1, 0], [0, 1]]
            while k:
                if k & 1:
                    res = multi(res, matrix)
                k >>= 1
                matrix = multi(matrix, matrix)
            return res

        res = qpow([[b, c], [1, 0]], n - 1)
        return res[0][0]


print(Solution().nthElement(2, 1, 1))
