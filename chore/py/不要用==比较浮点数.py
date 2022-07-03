import math


def check(assertion: str) -> None:
    print(assertion, end=" : ")
    print(eval(assertion))


check("1e50 == 10 ** 50")
check("abs(1e50 - 10 ** 50) < 1e-10")

# 64位最多表示浮点数 1e308 超出则为inf
# !对于inf和NaN会传染
check("1e500 == 1e600")
check("1e500 > 10 ** 1000")
check("1e500 * 1e500 > 0")
check("1e500 / 1e500 > 0")
check("1e500 / 1e500 == 1e500 / 1e500")


print(math.isinf(1e308), math.isinf(1e309))
print(math.isnan(1e500 / 1e500))


# 在有高安全的场景尽量避开浮点数 使用定点数
# 使用Decimal类 Decimal 数字的表示是完全精确的
