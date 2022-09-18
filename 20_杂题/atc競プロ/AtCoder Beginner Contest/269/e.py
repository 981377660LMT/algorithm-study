# 交互题

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def query(row1: int, row2: int, col1: int, col2: int) -> int:
    row1, row2, col1, col2 = row1 + 1, row2 + 1, col1 + 1, col2 + 1
    print(f"? {row1} {row2} {col1} {col2}", flush=True)
    return int(input())


def output(row: int, col: int) -> None:
    row, col = row + 1, col + 1
    print(f"! {row} {col}", flush=True)


# 寻找能够放置rook的位置
if __name__ == "__main__":
    n = int(input())
    # 二分法 先定行 再定列

    left1, right1 = 0, n - 1
    while left1 <= right1:
        mid = (left1 + right1) // 2
        count = query(0, mid, 0, n - 1)
        if count < mid + 1:
            right1 = mid - 1
        else:
            left1 = mid + 1

    left2, right2 = 0, n - 1
    while left2 <= right2:
        mid = (left2 + right2) // 2
        count = query(0, n - 1, 0, mid)
        if count < mid + 1:
            right2 = mid - 1
        else:
            left2 = mid + 1

    output(left1, left2)
