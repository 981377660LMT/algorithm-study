# https://atcoder.jp/contests/abc289/tasks/abc289_f
# 行列分开移动


# 二维平面，要求从点s到点 t。
# 给定一个矩形区域，每次操作从区域里选择一个点m，然后点 s就跑到与 点m的镜像点。
# 给定一个操作序列，或告知不可能。不要求最小化操作次数。

# 考虑一维
# 如果区域是一个点:s只有两种位置
# !如果区域不是一个点而是[a,b],那么两次操作分别选择a和a+1 就可以在这个方向移动2格.
# 依靠此基本操作一步一步移向终点。由坐标范围限制，操作次数不会超过上限。


from typing import List, Optional, Tuple


def teleportTakanashi(
    sx: int, sy: int, tx: int, ty: int, a: int, b: int, c: int, d: int
) -> Optional[List[Tuple[int, int]]]:
    def move(x: int, y: int, cx: int, cy: int) -> Tuple[int, int]:
        """起点(x,y)以(cx,cy)为中心,对称移动,并记录移动路径"""
        path.append((cx, cy))
        return cx - (x - cx), cy - (y - cy)

    path = []
    if (sx & 1) ^ (tx & 1) or (sy & 1) ^ (ty & 1):
        return
    if a == b and sx != tx:
        sx, sy = move(sx, sy, a, c)
    if c == d and sy != ty:
        sx, sy = move(sx, sy, a, c)
    if (a == b and sx != tx) or (c == d and sy != ty):
        return

    # 必定可以移动到终点
    while sx < tx:
        sx, sy = move(sx, sy, a, c)
        sx, sy = move(sx, sy, a + 1, c)
    while sx > tx:
        sx, sy = move(sx, sy, a + 1, c)
        sx, sy = move(sx, sy, a, c)
    while sy < ty:
        sx, sy = move(sx, sy, a, c)
        sx, sy = move(sx, sy, a, c + 1)
    while sy > ty:
        sx, sy = move(sx, sy, a, c + 1)
        sx, sy = move(sx, sy, a, c)

    return path


if __name__ == "__main__":
    sx, sy = map(int, input().split())
    tx, ty = map(int, input().split())
    a, b, c, d = map(int, input().split())
    path = teleportTakanashi(sx, sy, tx, ty, a, b, c, d)
    if path is None:
        print("No")
        exit(0)

    print("Yes")
    for a, b in path:
        print(a, b)
