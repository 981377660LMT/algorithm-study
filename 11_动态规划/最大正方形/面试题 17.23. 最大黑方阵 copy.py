from typing import List
from pprint import pprint


class Solution(object):
    def calRMatrix(self, matrix, rMatrix):
        for row in range(len(matrix) - 1, -1, -1):
            count = 0

            for col in range(len(matrix[0]) - 1, -1, -1):
                if matrix[row][col] == 1:
                    count = 0
                else:
                    count += 1

                rMatrix[row][col] = count

        return

    def calCMatrix(self, matrix, cMatrix):
        for col in range(len(matrix[0]) - 1, -1, -1):
            count = 0

            for row in range(len(matrix) - 1, -1, -1):
                if matrix[row][col] == 1:
                    count = 0
                else:
                    count += 1

                cMatrix[row][col] = count

        return

    def checkSquareValid(self, rMatrix, cMatrix, row, col, size):
        print(row, col, size)
        if row + size - 1 >= len(rMatrix) or col + size - 1 >= len(rMatrix[0]):
            return False

        rightUp = cMatrix[row][col + size - 1]
        leftDonw = rMatrix[row + size - 1][col]
        print(rightUp, leftDonw)
        if rightUp >= size and leftDonw >= size:
            return True

        return False

    def findSquare(self, matrix):
        """
        :type matrix: List[List[int]]
        :rtype: List[int]
        """
        rMatrix = [[0] * len(matrix[0]) for _ in range(len(matrix))]
        cMatrix = [[0] * len(matrix[0]) for _ in range(len(matrix))]
        resPos = [-1, -1]
        resSize = 0

        self.calRMatrix(matrix, rMatrix)
        self.calCMatrix(matrix, cMatrix)

        for row in range(len(matrix)):
            for col in range(len(matrix[0])):
                maxSize = min(rMatrix[row][col], cMatrix[row][col])

                while maxSize > 0:
                    if self.checkSquareValid(rMatrix, cMatrix, row, col, maxSize):
                        break

                    maxSize -= 1

                if maxSize > resSize:
                    resSize = maxSize
                    resPos = [row, col]
        pprint(rMatrix)
        pprint(cMatrix)
        if resSize == 0:
            return []
        else:
            return resPos + [resSize]


print(Solution().findSquare([[1, 0, 1], [0, 0, 1], [0, 0, 1],]))

