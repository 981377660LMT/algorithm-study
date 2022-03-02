from itertools import repeat, cycle, count, chain


print(*chain.from_iterable([[[1, 2, 3], [4, 5, 6]]]))


def linear(x, k, b=0):
    return k * x + b


a = {'x': 2, 'k': 3}
# 6，相当于 print(linear(x=2, k=3))
print(linear(**a))

