# 只支持加减乘括号的表达式解析，不能有空格

from typing import Tuple


def number(string: str, pos: int) -> Tuple[int, int]:
    res = 0
    while pos < len(string) and string[pos].isdigit():
        res = 10 * res + int(string[pos])
        pos += 1
    return res, pos


def factor(string: str, pos: int) -> Tuple[int, int]:
    if string[pos].isdigit():
        return number(string, pos)
    assert string[pos] == "("
    pos += 1
    res, pos = expr(string, pos)
    assert string[pos] == ")"
    pos += 1
    return res, pos


def term(string: str, pos: int) -> Tuple[int, int]:
    res, pos = factor(string, pos)
    while pos < len(string) and string[pos] == "*":
        pos += 1
        y, pos = factor(string, pos)
        res *= y
    return res, pos


def expr(string: str, pos: int) -> Tuple[int, int]:
    x, pos = term(string, pos)
    while pos < len(string):
        op = string[pos]
        if op != "+" and op != "-":
            return x, pos
        pos += 1
        y, pos = term(string, pos)
        if op == "+":
            x += y
        elif op == "-":
            x -= y
    return x, pos


def parse(string: str) -> int:
    return expr(string, 0)[0]


print(parse("1+2*3+(4*5+6)*7"))
