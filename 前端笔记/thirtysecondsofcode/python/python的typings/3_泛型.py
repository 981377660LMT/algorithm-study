from typing import Any, Callable, Generic, List, Mapping, TypeVar

X = TypeVar('X')
Y = TypeVar('Y')


def lookup_name(mapping: Mapping[X, Y], key: X, default: Y) -> Y:
    try:
        return mapping[key]
    except KeyError:
        return default


# lookup_name()

###########################################################
T = TypeVar('T')


class Stack(Generic[T]):
    def __init__(self) -> None:
        # Create an empty list with items of type T
        self.items: List[T] = []

    def push(self, item: T) -> None:
        self.items.append(item)

    def pop(self) -> T:
        return self.items.pop()

    def empty(self) -> bool:
        return not self.items


# 类型推断也适用于用户定义的泛型类型:
stack = Stack[int]()
stack.push(2)
stack.pop()
stack.push('x')  # Type error
###########################################################
# 泛型约束
StrOrInt = TypeVar('StrOrInt', str, int)

T = TypeVar('T')


class StrangePair(Generic[T, StrOrInt]):
    ...


###########################################################
# 泛型默认参数 bound
Func = TypeVar('Func', bound=Callable[..., Any])


def bare_decorator(func: Func) -> Func:
    ...


def decorator_args(url: str) -> Callable[[Func], Func]:
    ...

