from typing import Any


class Base:
    def foo(self, arg: Any) -> None:
        print("Base foo", arg)


class People(Base):
    ...


# 有100个继承People的类，有的实现了foo方法有的没有，且所有重写的foo方法只接受字符串类型参数
class Girl(People):
    ...


class Boy(People):
    def foo(self, arg: str) -> None:
        if not isinstance(arg, str):  # type: ignore
            raise TypeError(f"arg must be str, but got {arg})")
        print("Boy foo", arg)


if __name__ == "__main__":
    girl = Girl()
    girl.foo("hello")
    girl.foo(123)
    boy = Boy()
    boy.foo("hello")
    boy.foo(123)  # 要求：如果发现foo传入数字，需要转成字符串；请对上述代码重构，让代码运行时不报错

# 答案：


# 初始化顺序是 Boy -> People -> Base
# People 里的 self.foo 可以拿到当前实例最终调用的foo
class People(Base):
    def __init__(self):
        super().__init__()

        origin_foo = self.foo

        def f(arg: Any) -> None:
            if not isinstance(arg, str):
                arg = str(arg)
            origin_foo(arg)

        self.foo = f
