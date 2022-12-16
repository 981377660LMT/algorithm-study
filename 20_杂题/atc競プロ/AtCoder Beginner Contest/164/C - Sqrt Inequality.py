# 根号比较
# 判断 sqrt(a)+sqrt(b) < sqrt(c)是否成立
# !注意浮点数精度问题直接用Decimal
# 也可以公式变形(两边平方)


from decimal import Decimal


def sqrtSum1(a: int, b: int, c: int) -> bool:
    return 4 * a * b < (c - a - b) * (c - a - b) and c - a - b > 0


def sqrtSum2(a: int, b: int, c: int) -> bool:
    da, db, dc = Decimal(a), Decimal(b), Decimal(c)
    return da.sqrt() + db.sqrt() < dc.sqrt()


a, b, c = map(int, input().split())
print("Yes" if sqrtSum1(a, b, c) else "No")
