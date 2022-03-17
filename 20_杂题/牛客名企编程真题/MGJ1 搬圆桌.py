# 现在有一张半径为r的圆桌，其中心位于(x,y)，现在他想把圆桌的中心移到(x1,y1)。
# 每次移动一步，都必须在圆桌边缘固定一个点然后将圆桌绕这个点旋转。问最少需要移动几步。
import math

while True:
    try:
        r, x1, y1, x2, y2 = list(map(int, input().split()))

        # 注意每次旋转圆心都前进了2*r 贪心走斜边即可
        dist = math.dist((x1, y1), (x2, y2))

        print(math.ceil(dist / (2 * r)))

    except EOFError:
        break

