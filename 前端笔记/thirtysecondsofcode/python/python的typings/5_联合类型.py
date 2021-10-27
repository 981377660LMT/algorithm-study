from typing import Union

# 联合类型之联合类型会被展平
Union[Union[int, str], float] == Union[int, str, float]
# 在 3.10 版更改: 联合类型现在可以写成 X | Y
# StrOrInt = str | int
# Alternative syntax for unions requires Python 3.10 or newer

# 通常需要使用 isinstance ()检查来首先将联合类型缩小到非联合类型
def f(x: Union[int, str]) -> None:
    x + 1  # Error: str + int is not valid
    if isinstance(x, int):
        # Here type of x is int.
        x + 1  # OK
    else:
        # Here type of x is str.
        x + 'a'  # OK


f(1)  # OK
f('x')  # OK
f(1.1)  # Error
