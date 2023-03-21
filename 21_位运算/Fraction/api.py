# python分数类Fraction (没有Decimal快)


# 获取最简化分数形式后的分子和分母
import fractions

# !创建Fraction推荐用法，可以避免除法的浮点数误差
f1 = fractions.Fraction(2, -100)
print(f1)  # 已经是最简分数形式了


# !浮点数转为Fraction的最佳实践是用字符串包裹浮点数
print(fractions.Fraction(str(0.8857097)))  # 无误差 8857097/10000000
print(fractions.Fraction(0.8857097))  # 1994440937439217/2251799813685248 不对
print(fractions.Fraction(0.8857097).limit_denominator())  # 871913/984423  不对
print(fractions.Fraction.from_float(0.7))  # 3152519739159347/4503599627370496 不对

# 避免误差实例
print(0.1 + 0.2)
print(fractions.Fraction(1, 10) + fractions.Fraction(2, 10))

# 限制分母最大值
print(fractions.Fraction(23, 100).limit_denominator(10))

# 获取分母和分子
print(f1.numerator)
print(f1.denominator)
print(f1.as_integer_ratio())


# Fraction转为浮点数
print(float(f1))  # -0.02

# !注意要在函数外初始化Fraction 因为Fraction本身初始化计算量比较大
# 把 DISCOUNT = fractions.Fraction(7, 10) 放到外面初始化，里面直接用 a*DISCOUNT 会快上很多
