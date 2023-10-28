from functools import lru_cache


def main():
    n = int(input())
    pre = [[0] * 9 for _ in range(9)]
    arr = [tuple(map(int, input().split())) for _ in range(8)]

    # 获得题目要求公式的答案
    def get(x1, y1, x2, y2):
        tot = pre[x2][y2] - pre[x2][y1 - 1] - pre[x1 - 1][y2] + pre[x1 - 1][y1 - 1] - X
        return tot * tot / n

    @lru_cache(None)
    def mDFS(x1, y1, x2, y2, k):
        if k == 1:
            return get(x1, y1, x2, y2)
        res = 1e9
        # 枚举区间 [x1, x2] 中的间隙
        for i in range(x1, x2):
            # 得到分割线上方的数值 + 下方继续分割
            A = get(x1, y1, i, y2) + mDFS(i + 1, y1, x2, y2, k - 1)
            # 得到分割线下方的数值 + 上方继续分割
            B = get(i + 1, y1, x2, y2) + mDFS(x1, y1, i, y2, k - 1)
            res = min(res, A, B)
        # 枚举区间 [y1, y2] 中的间隙
        for j in range(y1, y2):
            # 同理，这里是分割为左右两块矩阵
            A = get(x1, y1, x2, j) + mDFS(x1, j + 1, x2, y2, k - 1)
            B = get(x1, j + 1, x2, y2) + mDFS(x1, y1, x2, j, k - 1)
            res = min(res, A, B)
        return res

    # 计算二维前缀和
    for i, line in enumerate(arr, 1):
        for j, x in enumerate(line, 1):
            pre[i][j] = pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1] + x

    # 计算 X 靶，最后补上开方以及保留小数位
    X = pre[8][8] / n
    print("%.3f" % mDFS(1, 1, 8, 8, n) ** 0.5)


if __name__ == "__main__":
    main()

