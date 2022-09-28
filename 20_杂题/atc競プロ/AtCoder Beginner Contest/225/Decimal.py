"""
Decimal

decimal意思为十进制,这个模块提供了十进制浮点运算支持
# !Decimal表示int 避免float运算的浮点数误差 (一般Decimal可以预处理)

可以传递给Decimal整型或者字符串参数,
`但不能是浮点数据`,因为浮点数据本身就不准确。

和Fraction模块的区别是什么?
Decimal比Fraction更快创建 (两倍左右)
Fraction(1,10) 可以写成 Decimal(1) / Decimal(10)
"""

from decimal import Decimal
from fractions import Fraction


a, b, c = Decimal(1), Decimal(2), Decimal(10)

print(a / c + b / c)  # 0.3
print(Decimal(6) * Decimal(3) / Decimal(9))
print(Fraction(6, 9))
