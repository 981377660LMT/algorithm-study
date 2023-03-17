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


def superKnight(n: int) -> int:
    """超级马经过n步后能抵达的格子数

    https://yukicoder.me/problems/no/1500
    二次多项式? => 待定系数法
    """
    if n <= 2:
        return [1, 12, 65][n]
    return (17 * n * n + 6 * n + 1) % 1000000007


print(superKnight(int(input())))
