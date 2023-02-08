# 距离√p最近的有理数表示

# https://tjkendev.github.io/procon-library/python/math/stern-brocot-tree.html


# a/b <= √p <= c/d を満たす a,b,c,d <= n を求める
from typing import Tuple


def stern_brocot(p: int, n: int) -> Tuple[int, int, int, int]:
    la = 0
    lb = 1
    ra = 1
    rb = 0
    lu = ru = 1
    lx = 0
    ly = 1
    rx = 1
    ry = 0
    while lu or ru:
        ma = la + ra
        mb = lb + rb
        if p * mb * mb < ma * ma:
            ra = ma
            rb = mb
            if ma <= n and mb <= n:
                rx = ma
                ry = mb
            else:
                lu = 0
        else:
            la = ma
            lb = mb
            if ma <= n and mb <= n:
                lx = ma
                ly = mb
            else:
                ru = 0

    # lx/ly <= √p <= rx/ry
    return lx, ly, rx, ry


print(stern_brocot(2, 10))
