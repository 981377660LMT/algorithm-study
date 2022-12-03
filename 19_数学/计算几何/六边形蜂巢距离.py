# 六边形距离公式(蜂巢距离公式)


def hexagon_distance(x1: int, y1: int, x2: int, y2: int) -> int:
    """蜂巢六边形中两点(x1, y1)和(x2, y2)的距离"""
    dx, dy = x1 - x2, y1 - y2
    return max(abs(dx), abs(dy), abs(dx - dy))  # 注意这也是三维空间中到原点的曼哈顿距离


assert hexagon_distance(0, 0, 1, 1) == 1
assert hexagon_distance(1, 1, 2, 0) == 2
