# https://www.cnblogs.com/Miracevin/p/9778266.html
# !三维前缀和


from typing import List


class MatrixPreSum3D:
    def __init__(self, matrix: List[List[List[int]]]) -> None:
        xSize, ySize, zSize = len(matrix), len(matrix[0]), len(matrix[0][0])
        preSum = [[[0] * (zSize + 1) for _ in range(ySize + 1)] for _ in range(xSize + 1)]
        for x in range(1, xSize + 1):
            for y in range(1, ySize + 1):
                for z in range(1, zSize + 1):
                    preSum[x][y][z] = matrix[x - 1][y - 1][z - 1]

        for x in range(1, xSize + 1):
            for y in range(1, ySize + 1):
                for z in range(1, zSize + 1):
                    preSum[x][y][z] += preSum[x - 1][y][z]

        for x in range(1, xSize + 1):
            for y in range(1, ySize + 1):
                for z in range(1, zSize + 1):
                    preSum[x][y][z] += preSum[x][y - 1][z]

        for x in range(1, xSize + 1):
            for y in range(1, ySize + 1):
                for z in range(1, zSize + 1):
                    preSum[x][y][z] += preSum[x][y][z - 1]

        self.preSum = preSum
        self.xSize, self.ySize, self.zSize = xSize, ySize, zSize

    def query(self, x1: int, y1: int, z1: int, x2: int, y2: int, z2: int) -> int:
        """
        查询 sum(A[x1:x2+1][y1:y2+1][z1:z2+1])的值

        MatrixPreSum3D.query(0, 0, 0, 1, 1, 1) == sum(A[0:2][0:2][0:2])
        """
        assert (
            0 <= x1 <= x2 < self.xSize and 0 <= y1 <= y2 < self.ySize and 0 <= z1 <= z2 < self.zSize
        )

        return (
            self.preSum[x2 + 1][y2 + 1][z2 + 1]
            - self.preSum[x1][y2 + 1][z2 + 1]
            - self.preSum[x2 + 1][y1][z2 + 1]
            - self.preSum[x2 + 1][y2 + 1][z1]
            + self.preSum[x1][y1][z2 + 1]
            + self.preSum[x1][y2 + 1][z1]
            + self.preSum[x2 + 1][y1][z1]
            - self.preSum[x1][y1][z1]
        )


if __name__ == "__main__":
    matrix3d = [
        [[1, 2, 3], [4, 5, 6], [7, 8, 9]],
        [[10, 11, 12], [13, 14, 15], [16, 17, 18]],
        [[19, 20, 21], [22, 23, 24], [25, 26, 27]],
    ]

    preSum3d = MatrixPreSum3D(matrix3d)
    assert preSum3d.query(1, 1, 1, 1, 1, 1) == 14
    assert preSum3d.query(0, 0, 0, 2, 2, 2) == sum(
        matrix3d[i][j][k] for i in range(3) for j in range(3) for k in range(3)
    )

    import time
    from numba import njit  # !numba加速for循环效果明显

    @njit
    def make3DArray(xSize: int, ySize: int, zSize: int) -> List[List[List[int]]]:
        res = [[[0] * zSize for _ in range(ySize)] for _ in range(xSize)]
        for x in range(xSize):
            for y in range(ySize):
                for z in range(zSize):
                    res[x][y][z] = x * y * z
        return res

    time1 = time.time()
    np3d = make3DArray(512, 512, 512)
    time2 = time.time()
    print(time2 - time1)

    # # 512*512*512
    # time2 = time.time()
    # np3d = np.array(matrix3d)
    # print(time.time() - time2)  # 7.779397487640381

    # time3 = time.time()
    # for i in range(1, 512):
    #     for j in range(1, 512):
    #         for k in range(1, 512):
    #             matrix3d[i][j][k] += matrix3d[i - 1][j][k]
    # print(time.time() - time3)  # 37.76209568977356
