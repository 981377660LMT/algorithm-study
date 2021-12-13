from typing import List


# 移动是指选择任一行或列，并转换该行或列中的每一个值：将所有 0 都更改为 1，将所有 1 都更改为 0。
# 在做出任意次数的移动后，将该矩阵的每一行都按照二进制数来解释，矩阵的得分就是这些数字的总和。
# 返回尽可能高的分数。

# 1.先将第 1 列为 0 的行反转成第 1 列全为 1
# 2.从第二列开始按列反转 如果当前列 1 的数量比 0 的数量要更少 则对当前列进行翻转


class Solution:
    def matrixScore(self, A: List[List[int]]) -> int:
        m, n = len(A), len(A[0])
        for i in range(m):
            if A[i][0] == 0:
                for j in range(n):
                    A[i][j] ^= 1

        for j in range(n):
            cnt = sum(A[i][j] for i in range(m))
            if cnt < m - cnt:
                for i in range(m):
                    A[i][j] ^= 1

        return sum(int("".join(map(str, A[i])), 2) for i in range(m))


print(Solution().matrixScore([[0, 0, 1, 1], [1, 0, 1, 0], [1, 1, 0, 0]]))
# 输出：39
# 解释：
# 转换为 [[1,1,1,1],[1,0,0,1],[1,1,1,1]]
# 0b1111 + 0b1001 + 0b1111 = 15 + 9 + 15 = 39
