def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
    if num < 0:
        return '-' + toString(-num, radix)

    if num == 0:
        return '0'

    res = []
    while num:
        div, mod = divmod(num, radix)
        res.append(str(mod))
        num = div
    return ''.join(res)[::-1]


def parseInt(string: str, radix: int) -> int:
    """将字符串转换为指定进制的数字"""
    return int(string, base=radix)

