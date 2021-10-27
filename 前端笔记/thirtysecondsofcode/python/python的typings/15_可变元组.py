from typing import Tuple


def f(t: Tuple[int, str]) -> None:
    t = 1, 'foo'  # OK
    # t = 'foo', 1  # Type check error


# as immutable, varying-length sequences
# 元组也可以用作不变的、可变长度的序列
def print_squared(t: Tuple[int, ...]) -> None:

    for n in t:
        print(n, n ** 2)


print_squared(())  # OK
print_squared((1, 3, 5))  # OK
print_squared([1, 2])  # Error: only a tuple is valid

# 通常使用 Sequence [ t ]代替 Tuple [ t，... ]是一个更好的主意，
# 因为 Sequence 也可以兼容 list 和其他非 Tuple 序列。
