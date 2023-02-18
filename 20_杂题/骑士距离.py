# 骑士距离
# 棋盘上马从(0,0)到目标点的最短路径


def knight_distance(tx: int, ty: int) -> int:
    """https://maspypy.github.io/library/other/knight_distance.hpp"""
    max = lambda x, y: x if x > y else y
    tx = abs(tx)
    ty = abs(ty)
    if tx + ty == 0:
        return 0
    if tx + ty == 1:
        return 3
    if tx == 2 and ty == 2:
        return 4
    step = (max(tx, ty) + 1) // 2
    step = max(step, (tx + ty + 2) // 3)
    step += (step ^ tx ^ ty) & 1
    return step


print(knight_distance(11, 9))
print(knight_distance(100, 100))
