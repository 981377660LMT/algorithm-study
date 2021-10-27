# Final names, methods and classes

# 最终名称是在初始化后不应该重新分配的变量或属性。它们对于声明常量很有用。
# 最终方法不应在子类中被重写(最后的方法)。
# 最后的类不应该被子类化(最后的类)。

from typing import Final, final

RATE: Final = 3000


class Base:
    DEFAULT_ID: Final = 0


RATE = 300  # Error: can't assign to final attribute
Base.DEFAULT_ID = 1  # Error: can't override a final attribute


# Final methods 最终方法
class Base:
    @final
    def common_name(self) -> None:
        ...


class Derived(Base):
    def common_name(self) -> None:
        ...  # Error: cannot override a final method


# Final classes


@final
class Leaf:
    ...


class MyLeaf(Leaf):  # Error: Leaf can't be subclassed
    ...
