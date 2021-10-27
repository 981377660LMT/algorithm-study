class Base:
    def f(self, x: int) -> None:
        ...


class Derived1(Base):
    def f(self, x: str) -> None:  # Error: type of 'x' incompatible
        ...


class Derived2(Base):
    def f(self, x: int, y: int) -> None:  # Error: too many arguments
        ...


class Derived3(Base):
    def f(self, x: int) -> None:  # OK
        ...


class Derived4(Base):
    def f(self, x: float) -> None:  # OK: mypy treats int as a subtype of float
        ...


class Derived5(Base):
    def f(self, x: int, y: int = 0) -> None:  # OK: accepts more than the base
        ...  # class method

