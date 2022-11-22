class MyTemplate:
    attr: str

    def foo(self) -> None:
        ...

    def bar(self) -> None:
        ...


myList = ["a", "b", "c", "d"]


def genClass(arg: str):
    return type("MyClass", (MyTemplate,), {"attr": arg})


cls1 = genClass("a")()
cls1.foo()
