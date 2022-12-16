# !使用Decimal避免浮点数乘法精度问题
# 求a*b 取整,其中a,b为浮点数
from decimal import Decimal


if __name__ == "__main__":
    A, B = input().split()
    num1, num2 = Decimal(A), Decimal(B)
    print(int(num1 * num2))
