a, b, c, d = 1, 1, 1000, 1000
print(a is b, c is d)


def foo():
    e = 1000
    f = 1000
    print(e is f, e is d)
    g = 1
    print(g is a)


foo()
