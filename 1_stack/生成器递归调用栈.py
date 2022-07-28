import sys
from types import GeneratorType
from typing import Any, Callable, Generator, TypeVar

sys.setrecursionlimit(int(1e6))


def fib1(n: int, a=0, b=1):
    if n <= 1:
        yield b
        return
    yield from fib1(n - 1, b, a + b)


# print(*fib1(1000))  # !RecursionError: maximum recursion depth exceeded in comparison
##############################################################################################
# https://atcoder.jp/contests/abc261/submissions/33481104
C = TypeVar("C", bound=Callable[..., Generator[Any, Any, Any]])


def bootstrap(func: C) -> C:

    stack = []

    def wrapper(*args, **kwargs) -> Generator[Any, Any, Any]:
        if stack:
            return func(*args, **kwargs)

        to = func(*args, **kwargs)
        while True:
            if type(to) is GeneratorType:
                stack.append(to)  # 生成器前序遍历
                to = next(to)
            else:
                stack.pop()
                if not stack:
                    break
                to = stack[-1].send(to)  # 生成器后序遍历

        return to

    return wrapper  # type: ignore


fib2 = bootstrap(fib1)


res = fib2(100)
print(res)
