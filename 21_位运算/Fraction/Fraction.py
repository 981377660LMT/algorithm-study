# https://kazun1998.github.io/library_for_python/Fraction.py

from math import gcd


class Fraction:
    reduction = False
    expand = False

    """
    reduction : 分数を約分した状態で保存するかどうか
    expand : 1/0, -1/0 を認めるかどうか
    """

    __slots__ = ("__a", "__b")

    def __init__(self, Numerator=0, Denominator=1):
        """分数のオブジェクトを生成する.

        Numerator : 分子
        Denominator : 分母
        """
        assert Denominator or Fraction.expand, "分母が0"

        if Denominator < 0:
            Numerator *= -1
            Denominator *= -1

        self.__a = Numerator
        self.__b = Denominator
        if Fraction.reduction:
            g = gcd(Numerator, Denominator)
            self.__a = Numerator // g
            self.__b = Denominator // g

    def numerator(self):
        return self.__a

    def denominator(self):
        return self.__b

    def __iter__(self):
        yield self.__a
        yield self.__b

    # 表示定義
    def __str__(self):
        if self.__b == 1:
            return str(self.__a)
        else:
            return f"{self.__a}/{self.__b}"

    __repr__ = __str__

    # 四則演算定義
    def __add__(self, other):
        if isinstance(other, Fraction):
            x = self.__a * other.__b + self.__b * other.__a
            y = self.__b * other.__b
        elif isinstance(other, int):
            x = self.__a + self.__b * other
            y = self.__b
        else:
            raise TypeError(f"{other} の型がおかしいです.")
        return Fraction(x, y)

    def __radd__(self, other):
        return self + other

    def __sub__(self, other):
        if isinstance(other, Fraction):
            x = self.__a * other.__b - self.__b * other.__a
            y = self.__b * other.__b
        elif isinstance(other, int):
            x = self.__a - self.__b * other
            y = self.__b
        else:
            raise TypeError(f"{other} の型がおかしいです.")
        return Fraction(x, y)

    def __rsub__(self, other):
        return -self + other

    def __mul__(self, other):
        if isinstance(other, Fraction):
            x = self.__a * other.__a
            y = self.__b * other.__b
        elif isinstance(other, int):
            x = self.__a * other
            y = self.__b
        else:
            raise TypeError(f"{other} の型がおかしいです.")

        return Fraction(x, y)

    def __rmul__(self, other):
        return self * other

    def __floordiv__(self, other):
        if other == Fraction():
            raise ZeroDivisionError

        H = self / other
        return H.a // H.b

    def __rfloordiv__(self, other):
        if self == Fraction():
            raise ZeroDivisionError

        H = other / self
        return H.a // H.b

    def __truediv__(self, other):
        if not other:
            raise ZeroDivisionError("ゼロ除算")

        if isinstance(other, Fraction):
            x = self.__a * other.__b
            y = self.__b * other.__a
        elif other.__class__ == int:
            x = self.__a
            y = self.__b * other
        else:
            raise TypeError(f"{other} の型がおかしいです.")

        return Fraction(x, y)

    def __rtruediv__(self, other):
        if not self:
            raise ZeroDivisionError("ゼロ除算")

        if isinstance(other, Fraction):
            x = other.__a * self.__b
            y = other.__b * self.__a
        elif isinstance(other, int):
            x = other * self.__b
            y = self.__a
        else:
            raise TypeError(f"{other} の型がおかしいです.")
        return Fraction(x, y)

    def __pow__(self, m):
        alpha, beta = self

        if m < 0:
            alpha, beta = beta, alpha
            m = -m

        return Fraction(pow(alpha, m), pow(beta, m))

    # 丸め
    def __floor__(self):
        return self.__a // self.__b

    def __ceil__(self):
        return (self.__a + self.__b - 1) // self.__b

    # 真偽値
    def __bool__(self):
        return bool(self.__a)

    def __compare(self, other):
        """self < other ならば -1, self = other ならば 0, self > other ならば, +1"""
        if isinstance(other, Fraction):
            x = self.__a * other.__b
            y = self.__b * other.__a
        else:
            x = self.__a
            y = self.__b * other

        if x < y:
            return -1
        elif x > y:
            return 1
        else:
            return 0

    # 比較
    def __eq__(self, other):
        return self.__compare(other) == 0

    def __neq__(self, other):
        return self.__compare(other) != 0

    def __lt__(self, other):
        return self.__compare(other) == -1

    def __le__(self, other):
        return self.__compare(other) != 1

    def __gt__(self, other):
        return self.__compare(other) == 1

    def __ge__(self, other):
        return self.__compare(other) != -1

    # その他
    def __float__(self):
        return self.__a / self.__b

    def sign(self):
        if self.__a > 0:
            return 1
        elif self.__a == 0:
            return 0
        else:
            return -1

    # 符号
    def __pos__(self):
        return self

    def __neg__(self):
        return Fraction(-self.__a, self.__b)

    def __abs__(self):
        if self.__a > 0:
            return self
        else:
            return -self

    # その他
    def is_unit(self):
        return self.__a == 1

    def __hash__(self):
        x, y = self
        if not Fraction.reduction:
            g = gcd(x, y)
            x //= g
            y //= g
        return hash((x, y))


if __name__ == "__main__":
    # 启用约分（可选，根据需要设置）
    Fraction.reduction = True
    Fraction.expand = False  # 不允许 1/0, -1/0

    # 创建几个 Fraction 对象
    f1 = Fraction(1, 2)
    f2 = Fraction(3, 4)
    f3 = Fraction(5, 6)
    f4 = Fraction(-2, 3)

    print("f1 =", f1)  # 输出 "1/2"
    print("f2 =", f2)  # 输出 "3/4"
    print("f3 =", f3)  # 输出 "5/6"
    print("f4 =", f4)  # 输出 "-2/3"

    # 四则运算
    print("f1 + f2 =", f1 + f2)  # 1/2 + 3/4 = 5/4
    print("f2 - f1 =", f2 - f1)  # 3/4 - 1/2 = 1/4
    print("f1 * f3 =", f1 * f3)  # 1/2 * 5/6 = 5/12
    print("f3 / f2 =", f3 / f2)  # 5/6 / 3/4 = 10/9

    # 幂运算
    print("f1 ** 2 =", f1**2)  # (1/2)^2 = 1/4
    print("f4 ** 3 =", f4**3)  # (-2/3)^3 = -8/27

    # 浮点数转换
    print("float(f1) =", float(f1))  # 0.5

    # 比较
    print("f1 < f2:", f1 < f2)  # True
    print("f1 == Fraction(2, 4):", f1 == Fraction(2, 4))  # True

    # 求绝对值
    print("abs(f4) =", abs(f4))  # 2/3

    # 使用迭代获得分子和分母
    a, b = f1
    print("f1 numerator =", a)
    print("f1 denominator =", b)
