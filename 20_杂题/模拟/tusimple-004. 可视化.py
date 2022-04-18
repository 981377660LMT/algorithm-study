# n<=100
# 输出无限大矩阵中一个包含所有有效信息的整行的最小矩阵。
# 如果净资产在第 i 天增加，则他将在某行中绘制 / ，其行索引等于当日开始时的净资产。
# 如果净资产在第 i 天减少，则他将在某行中绘制 \ ，其行索引等于当日结束时的净资产。
# 如果第 i 天的净值未发生变化，则他将在某行中绘制 - ，其行索引等于该天的净资产 。
# 所有其他单元格都用 . 填充。
# 其实应该记录每天资产两个端点，然后每取天最小值作为行索引，记得要用最大值平移


# 绘图题
n = int(input())
s = input()

mapping = {'+': '/', '-': '\\', '=': '-'}


# 直接一遍遍历，字典存储每个点的高度


preChar, preRow = s[0], 0  # 之前的字符，之前的行编号，越往上负的越多
upRow, BottomRow = 0, 0
for col, char in enumerate(s[1:], start=1):
    curRow = -1
    if char == '+':
        curRow = preRow - (1 if preChar == '+' else 0)
    elif char == '-':
        curRow = preRow + (1 if preChar == '-' else 0)
    else:
        curRow = preRow

    upRow = min(upRow, curRow)
    BottomRow = max(BottomRow, curRow)
    preChar, preRow = char, curRow


ROW, COL = BottomRow - upRow + 1, n
matrix = [['.'] * COL for _ in range(ROW)]


# row 平移
preChar, preRow = s[0], -upRow
matrix[preRow][0] = mapping[preChar]
for col, char in enumerate(s[1:], start=1):
    curRow = -1
    if char == '+':
        curRow = preRow - (1 if preChar == '+' else 0)
    elif char == '-':
        curRow = preRow + (1 if preChar == '-' else 0)
    else:
        curRow = preRow

    matrix[curRow][col] = mapping[char]
    preChar, preRow = char, curRow


for row in matrix:
    print(''.join(row))


# +++ 不对
# ++- 不对
