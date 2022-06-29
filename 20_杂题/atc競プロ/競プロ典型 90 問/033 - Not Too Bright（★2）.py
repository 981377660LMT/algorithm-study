# 矩阵中需要放置LED灯
# 注意LED灯为左上角2x2的方格中不能有其他LED灯
# 求放置LED灯的最大数量


ROW, COL = map(int, input().split())
if ROW == 1 or COL == 1:
    print(ROW * COL)
else:
    print(((ROW + 1) // 2) * ((COL + 1) // 2))  # 注意括号

