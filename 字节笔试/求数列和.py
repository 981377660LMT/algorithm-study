# 数列的第一项为n，以后各项为前一项的平方根，求数列的前m项的和。

from math import sqrt


while True:
    try:
        n, m = map(int, input().split())
        cur = n
        res = 0
        for _ in range(m):
            res += cur
            cur = sqrt(cur)
        # 要求精度保留2位小数。
        print(format(res, '.2f'))
    except EOFError:
        break

