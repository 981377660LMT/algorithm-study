多点共线

# 计算斜率，注意分母不为 0:

1. Fraction((y2 - y1) / (x2 - x1)))
2. 元组
   gcd* = gcd(x2 - x1, y2 - y1)
   ((y2 - y1) // gcd*, (x2 - x1) // gcd\_)

```Python
计算斜率key，key用元组表示；gcd最好自己写

A, B = y2 - y1, x2 - x1
if B == 0:
    key = (0, 0)
else:
    gcd_ = gcd(A, B)
    key = (A / gcd_, B / gcd_)
```

多点共线问题参考
`面试题 16.14. 最佳直线 copy`
`2152.覆盖点所需要的最少直线数`
