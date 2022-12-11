# C - Super Ryuma-超龙马
# 超龙马可以走到距离当前格子曼哈顿距离为3以内的格子
# 也可以走到四条对角线上的格子
# 问从起点到终点有多少种走法

# !答案在0,1,2,3中
def superRyuma(sr: int, sc: int, tr: int, tc: int) -> int:
    rowDiff, colDiff = tr - sr, tc - sc
    if rowDiff == 0 and colDiff == 0:
        return 0
    if rowDiff == colDiff or rowDiff == -colDiff:  # 对角线
        return 1
    if abs(rowDiff) + abs(colDiff) <= 3:  # 曼哈顿距离为3以内
        return 1
    if not (rowDiff ^ colDiff) & 1:  # 同色
        return 2
    if abs(rowDiff) + abs(colDiff) <= 6:  # 两次曼哈顿距离为6以内
        return 2
    if abs(rowDiff + colDiff) <= 3 or abs(rowDiff - colDiff) <= 3:  # 斜对角+曼哈顿距离3以内
        return 2
    return 3


if __name__ == "__main__":
    sr, sc = map(int, input().split())
    tr, tc = map(int, input().split())
    print(superRyuma(sr, sc, tr, tc))
