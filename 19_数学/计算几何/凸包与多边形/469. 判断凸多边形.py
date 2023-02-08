# 给定一个按顺序连接的多边形的顶点，判断该多边形是否为凸多边形
from typing import List

# 叉乘判断
# 设A(x1,y1),B(x2,y2),C(x3,y3)则三角形两边的矢量分别是：
# AB=(x2-x1,y2-y1), AC=(x3-x1,y3-y1)
# 则AB和AC的叉积为：(2*2的行列式) 值为：(x2-x1)*(y3-y1) - (y2-y1)*(x3-x1)
# 利用右手法则进行判断：
# 如果AB*AC>0,则三角形ABC是逆时针的
# 如果AB*AC<0,则三角形ABC是顺时针的


def isConvex(points: List[List[int]]) -> bool:
    def cross(A, B, C):
        AB = [B[0] - A[0], B[1] - A[1]]
        AC = [C[0] - A[0], C[1] - A[1]]
        return AB[0] * AC[1] - AB[1] * AC[0]

    flag = 0
    n = len(points)
    for i in range(n):
        # cur > 0 表示points是按逆时针输出的;cur < 0,顺时针
        cur = cross(points[i], points[(i + 1) % n], points[(i + 2) % n])
        if cur != 0:
            # 说明异号, 说明有个角大于180度
            if cur * flag < 0:
                return False
            else:
                flag = cur
    return True


# 两条线段：利用点相减的形式表示，则AC = C - A =（0, 2），AB = B - A =（1, 1）
# 线段a与线段b的叉积，记作a×b，公式：若a =（x1,y1），b= (x2,y2)，则
# AC×AB = x1y2 - x2y1
# 所以样例中，AB×AC=1×2-0×1=1
# 扯了半天，我们得出：AB×AC的结果为1，大于0
# 根据结果的正负关系，我们可以推出AB与AC的关系（这里的顺逆时针是指B、C点围绕极点A的顺逆关系）：
# ①若结果大于0，则AC在AB的逆时针方向(a->b->c 逆时针)
# ②若结果小于0，则AC在AB的顺时针方向(a->b->c 顺时针)
# ③若结果等于0，则AC与AB平行（其实是共线，但是可能朝向相反）
# 或者我们验证一下这个结论，将AB×AC改为AC×AB，则得到结果：
# AC×AB = 0×1-1×2 = -1，结果刚好相反，可以理解为交换了计算顺序，即是互换了线段位置，也符合结论。

# 作者：fly-f
# 链接：https://leetcode-cn.com/problems/erect-the-fence/solution/xiang-xi-ti-jie-by-fly-f/
# 来源：力扣（LeetCode）
# 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
