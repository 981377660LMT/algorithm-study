# 求关于 x 的同余方程 ax≡1(modb) 的最小正整数解。
from typing import Tuple


a, b = map(int, input().split())

# python 3.5 版本报错：pow() 2nd argument cannot be negative when 3rd argument specified
# print(pow(a, -1, b))

# 欧几里得算法原理：(a,b)=(b,a%b)
# 扩展欧几里得求逆元


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """
    求a, b最大公约数，同时求出裴蜀定理中的一组系数x, y， 满足x*a + y*b = gcd(a, b)
    ax+by=gcd_ 返回 (gcd_,x,y)
    """
    if b == 0:
        return a, 1, 0
    else:
        gcd_, x, y = exgcd(b, a % b)
        return gcd_, y, x - a // b * y


gcd_, x, y = exgcd(a, b)  # ax+by=gcd_
print(x % b)

