def debug(fun, DEBUG=float('inf')):
    """
    打印函数的所有输入与输出.
    可以通过DEBUG设置递归深度
    """

    def deco(*args, **kw):
        nonlocal DEBUG
        DEBUG -= 1
        result = fun(*args, **kw)
        p, t, block = args
        if DEBUG > 0:
            print(fun.__name__, args, kw if kw else '', '=', result)
        DEBUG += 1
        return result

    return deco
