# rationalGcd  有理数gcd (有理数通分)


import math
from fractions import Fraction


def qlcm(f1: "Fraction", f2: "Fraction") -> "Fraction":
    a = f1.numerator * f2.denominator
    b = f1.denominator * f2.numerator
    c = f1.denominator * f2.denominator
    return Fraction(math.lcm(a, b), c)


def qgcd(f1: "Fraction", f2: "Fraction") -> "Fraction":
    return f1 * f2 / qlcm(f1, f2)


if __name__ == "__main__":
    f1 = Fraction(1, 2)
    f2 = Fraction(1, 3)
    print(qgcd(f1, f2))
