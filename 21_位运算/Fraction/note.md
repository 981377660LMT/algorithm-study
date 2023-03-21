在 Python 中，Decimal 和 Fraction 是两种用于表示非整数数值的数据类型，它们都可以提供比浮点数更高的精度。然而，它们的表示方式和适用场景有所不同。

## Decimal

Decimal 类来自 decimal 模块，它用于表示十进制浮点数。Decimal 类在精度和范围方面提供了更好的控制，并遵循 IEEE 754 标准。Decimal 适用于需要精确计算的场景，例如金融和货币计算。

```Python
from decimal import Decimal

a = Decimal('0.1')
b = Decimal('0.2')
c = a + b

print(c)  # 输出：0.3


```

## Fraction

Fraction 类来自 fractions 模块，它用于表示有理数，即分数。Fraction 类通过两个整数（分子和分母）表示一个有理数。Fraction 在需要进行精确数学运算的场景中很有用，尤其是在涉及分数计算的场景中。

```python
from fractions import Fraction

a = Fraction(1, 3)
b = Fraction(1, 6)
c = a + b

print(c)  # 输出：1/2

```

## 总结：

Decimal 用于表示`十进制浮点数`，适用于精确计算，尤其是金融和货币计算。
Fraction 用于表示`有理数（分数`），适用于精确数学运算，尤其是涉及分数计算的场景。
