# 在一个 个方格组成的棋盘中，有一个方格与其它的不同，使用四种L型骨牌覆盖除这个特殊方格的其它方格，请使用分治法实现棋盘覆盖
def chess(startRow: int, startCol: int, specialRow: int, specialCol: int, size: int) -> None:
    # tr:棋盘初始行号 tc:棋盘初始列号
    # pr:特殊棋盘格行号 pc:特殊棋盘格列号
    # size:棋盘格大小
    global mark
    global table
    if size == 1:
        return  # 递归终止条件
    mark += 1  # 表示直角骨牌号
    count = mark
    half = size // 2  # 当size不等于1时，棋盘格规模减半，变为4个

    # 小棋盘格进行递归操作
    # 左上角
    if (specialRow < startRow + half) and (specialCol < startCol + half):
        chess(startRow, startCol, specialRow, specialCol, half)
    else:
        table[startRow + half - 1][startCol + half - 1] = count
        chess(startRow, startCol, startRow + half - 1, startCol + half - 1, half)
        # 将[tr+half-1,tc+half-1]作为小规模棋盘格的特殊点，进行递归

    # 右上角
    if (specialRow < startRow + half) and (specialCol >= startCol + half):
        chess(startRow, startCol + half, specialRow, specialCol, half)
    else:
        table[startRow + half - 1][startCol + half] = count
        chess(startRow, startCol + half, startRow + half - 1, startCol + half, half)
    # 将[tr+half-1,tc+half]作为小规模棋盘格的特殊点，进行递归

    # 左下角
    if (specialRow >= startRow + half) and (specialCol < startCol + half):
        chess(startRow + half, startCol, specialRow, specialCol, half)
    else:
        table[startRow + half][startCol + half - 1] = count
        chess(startRow + half, startCol, startRow + half, startCol + half - 1, half)
    # 将[tr+half,tc+half-1]作为小规模棋盘格的特殊点，进行递归

    # 右下角
    if (specialRow >= startRow + half) and (specialCol >= startCol + half):
        chess(startRow + half, startCol + half, specialRow, specialCol, half)
    else:
        table[startRow + half][startCol + half] = count
        chess(startRow + half, startCol + half, startRow + half, startCol + half, half)
    # 将[tr+half,tc+half]作为小规模棋盘格的特殊点，进行递归


def show(table):
    n = len(table)
    for i in range(n):
        for j in range(n):
            print(table[i][j], end='	')
        print('')


mark = 0
n = 8  # 输入8*8的棋盘规格
table = [[-1 for x in range(n)] for y in range(n)]  # -1代表特殊格子
chess(0, 0, 2, 2, n)  # 特殊棋盘位置
show(table)
