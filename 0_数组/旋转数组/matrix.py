# 2次元行列を右に90度回転
def matrix_rotate(A):
    H, W = len(A), len(A[0])
    ans = [[None] * H for _ in range(W)]

    for y in range(H):
        for x in range(W):
            ans[W - 1 - x][y] = A[y][x]

    return ans


## matrix_shift


# 2 次元行列 A の y 行を右にシフト
def shift_right_y(y, A):
    back = A[y][-1]
    for x in range(len(A[0]) - 1, 0, -1):
        A[y][x] = A[y][x - 1]
    A[y][0] = back


# 2 次元行列 A の y 行を左にシフト
def shift_left_y(y, A):
    back = A[y][0]
    for x in range(1, len(A[0])):
        A[y][x - 1] = A[y][x]
    A[y][-1] = back


# 2 次元行列 A の x 列を上にシフト
def shift_up_x(x, A):
    back = A[0][x]
    for y in range(1, len(A)):
        A[y - 1][x] = A[y][x]
    A[-1][x] = back


# 2 次元行列 A の x 列を下にシフト
def shift_down_x(x, A):
    back = A[-1][x]
    for y in range(len(A) - 1, 0, -1):
        A[y][x] = A[y - 1][x]
    A[0][x] = back
