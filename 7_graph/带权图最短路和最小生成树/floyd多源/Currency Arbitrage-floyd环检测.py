# 货币汇率转换
# n<=100 应该n^3


class Solution:
    def solve(self, matrix):
        # 判断是否能钱生钱,即存在环
        n = len(matrix)
        for i in range(n):
            for j in range(n):
                for k in range(n):
                    matrix[i][j] = max(matrix[i][k] * matrix[k][j], matrix[i][j])
            if matrix[i][i] > 1:
                return True
        return False


print(Solution().solve(matrix=[[1, 1.32, 0.9], [0.76, 1, 0.72], [1.11, 1.47, 1]]))

# The value at entry [i][j] in this matrix represents the amount of currency j you could buy with one unit of currency i. Let's say currency 0 is USD, 1 is CAD and 2 is EUR. We can make an arbitrage with the following:

# Sell 1 CAD for 0.72 EUR
# Sell 0.72 EUR for 0.80 USD (0.72 * 1.11)
# Sell 0.80 USD for 1.055 CAD (0.72 * 1.11 * 1.32)
