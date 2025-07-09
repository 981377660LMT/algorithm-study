from math import gcd
from typing import List, Tuple


# 给你一个无穷大的网格图。一开始你在 (1, 1) ，你需要通过有限步移动到达点 (targetX, targetY) 。
# 每一步 ，你可以从点 (x, y) 移动到以下点之一：
# (x, y - x)
# (x - y, y)
# (2 * x, y)
# (x, 2 * y)
# 给你两个整数 targetX 和 targetY ，分别表示你最后需要到达点的 X 和 Y 坐标。
# !如果你可以从 (1, 1) 出发到达这个点，请你返回true ，否则返回 false 。


# Binary Gcd(二进制gcd)
# 解:
# 两个正向操作是 更相减损术 的步骤，出现在高中数学必修三里，由这个可以往 gcd 联想
# 反过来思考，操作变为(x,x+y),(x+y,y),(x/2,y),(x,y/2)。
# 如果两个数存在2以外的公约数，显然所有操作都不能约掉这个公约数。
# 如果不存在这样的公约数（除2以外互质），我们总是可以通过这样的操作将两个数之和持续变小：
# !1. 如果两个数中有偶数，将偶数除以2
# !2. 如果两个数都是奇数，保留较小的数，较大的数加上较小的数（必然是偶数）然后立即除以2，
#    相当于（假设x<y)，变成(x, (x+y)/2)
#    以上操作总是将x+y变得更小，直到x=y，又因为初始的xy没有2以外的公约数，所以显然x=y=1


class Solution:
    def isReachable(self, targetX: int, targetY: int) -> bool:
        return gcd(targetX, targetY).bit_count() == 1  # gcd是2的幂次
        while targetX % 2 == 0:
            targetX //= 2
        while targetY % 2 == 0:
            targetY //= 2
        return gcd(targetX, targetY) == 1

    def getPath(self, targetX: int, targetY: int) -> List[Tuple[int, int]]:
        """逆向思维"""
        if not self.isReachable(targetX, targetY):
            return []

        path = []
        while (targetX, targetY) != (1, 1):
            path.append((targetX, targetY))
            if targetX % 2 == 0:
                targetX //= 2
                continue
            if targetY % 2 == 0:
                targetY //= 2
                continue
            if targetX > targetY:
                targetX += targetY
            else:
                targetY += targetX

        path.append((1, 1))
        return path[::-1]


print(Solution().getPath(targetX=2, targetY=3))
print(Solution().getPath(targetX=4, targetY=7))


def binaryGcd(a: int, b: int) -> int:
    """
    gcd(a,b) =
    - a if a == b
    - gcd(a//2,b//2) if a,b are even
    - gcd(a//2,b) if a is even and b is odd
    - gcd((a-b)//2,b) if a,b are odd


    https://zhuanlan.zhihu.com/p/553890800
    https://riteme.site/blog/2016-8-19/binary-gcd.html
    """
    if a < b:  # 保证 a >= b
        a, b = b, a
    if b == 0:
        return a

    if a & 1:
        if b & 1:
            return binaryGcd((a - b) >> 1, b)
        return binaryGcd(a, b >> 1)

    if b & 1:
        return binaryGcd(a >> 1, b)
    return binaryGcd(a >> 1, b >> 1) << 1
