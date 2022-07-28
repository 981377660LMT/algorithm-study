# https://docs.python.org/3/library/typing.html
# https://github.com/microsoft/python-type-stubs
# https://google.github.io/pytype/developers/type_stubs.html

import sys
from typing import (
    Annotated,
    Any,
    Callable,
    ClassVar,
    Final,
    NewType,
    NoReturn,
    Optional,
    ParamSpec,
    Protocol,
    Type,
    TypeVar,
    cast,
    final,
    runtime_checkable,
    TypedDict,
)
import typing_extensions

# from typing_extensions import LiteralString


# !1.类变量(static)和实例变量定义 用ClassVar来区分
class Foo:
    clsValue: ClassVar[int] = 0
    instValue: int = 0


foo = Foo()
foo.clsValue = (
    1  # Member "clsValue" cannot be assigned through a class instance because it is a ClassVar
)
foo.instValue = 1

# !2. NewType来严格区分不同的int类型
UserId = NewType("UserId", int)
StaffId = NewType("StaffId", int)
id1 = UserId(1)
id2 = StaffId(1)
a: UserId = 1  # "Literal[1]" is incompatible with "UserId"
print(id1 == id2)

# !3. Any 和 object 的区别 尽量使用object
# object是所有类型的父类，反则不成立
def hash1(item: object) -> int:
    return item.magic()  # 类似ts里的unknown类型


def hash2(item: Any) -> int:
    return item.magic()


# !4. 使得protocol接口变为运行时检查 (很多protocol都以Supportsxxx命名)
@runtime_checkable
class SupportsBar(Protocol):
    __slots__ = ()

    def bar(self) -> NoReturn:  # never
        raise NotImplementedError("bar")


# !5. 函数的泛型：ParamSpec 保留装饰器函数类型信息  Concatenate
# https://sobolevn.me/2021/12/paramspec-guide
# 它们分别采取 Callable[ParamSpecVariable, ReturnType]
# 和 Callable[Concatenate[Arg1Type, Arg2Type, ..., ParamSpecVariable], ReturnType] 的形式。
#  Concatenate 目前只在作为 Callable 的第一个参数时有效。Concatenate 的最后一个参数必须是一个 ParamSpec。
# ParamSpec 里又有 ParamSpecArgs 和 ParamSpecKwargs 分别表示参数的类型和名称。

# !原来的装饰器：
C = TypeVar("C", bound=Callable[..., Any])


def logger(func: C) -> C:
    def wrapper(*args: object, **kwargs: object):
        print("logging")
        return func(*args, **kwargs)

    return wrapper  # type: ignore


# !现在的装饰器


P = ParamSpec("P")  # 类似ts的ParameterType
R = TypeVar("R")


def addLogging(f: Callable[P, R]) -> Callable[P, R]:
    def wrapper(*args: object, **kwargs: object) -> R:
        print("logging")
        return f(*args, **kwargs)  # type: ignore

    return wrapper  # type: ignore


@addLogging
def add(x: int, y: int) -> int:
    return x + y


# !如果想要连接函数参数，可以使用 Concatenate。

# !6. Type获取实例变量类型(类似ts的typeof)
foo = Foo()


def modifyFoo(foo: Type[Foo]) -> None:
    foo.instValue = 1


# !7. Final/@final 告知类型检查器某名称不能再次赋值或在子类中重写的特殊类型构造器 类似readonly
@final  # !cannot be overridden, and decorated class cannot be subclassed
class FinalClass:
    abc: Final[int] = 1


class SubClass(FinalClass):
    def __init__(self):
        self.abc = 2  # cannot be redeclared because parent class "FinalClass" declares it as Final


# !8. 类型断言cast函数(类似id函数)
a: str | int = 1
a = cast(str, a)
print(a.capitalize())


# !9. 限制字典的key类型 使用 TypedDict
class SDict(TypedDict):
    a: int
    b: str
    c: Optional[str]


sd: SDict = {"a": 1, "c": "3", "b": "2"}

# !新版特性在extension_typings里:Unpack解构元素类型 todo/Self/Never/LiteralString

print(sys.version_info >= (3, 11))
print(sys.version_info)  # 命名元组
