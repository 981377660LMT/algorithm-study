import numpy as np


# !多项式乘法(处理生成函数)
poly1 = np.poly1d([1, 2, 3])  # x^2 + 2x + 3
poly2 = np.poly1d([4, 5, 6])  # 4x^2 + 5x + 6
print(poly1 * poly2)  # 4x^4 + 13x^3 + 28x^2 + 27x + 18
print(poly1.coef.tolist())  # [1, 2, 3]

# !多项式除法
poly3 = poly2 / poly1
div, mod = poly3  # 分别为商式和余式
print(div, mod)
print(div.coef.tolist())  # [4.0]
print(mod.coef.tolist())  # [-3.0, -6.0]
