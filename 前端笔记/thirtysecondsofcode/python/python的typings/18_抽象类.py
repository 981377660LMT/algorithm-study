from abc import ABCMeta, abstractmethod


class Animal(metaclass=ABCMeta):
    @abstractmethod
    def eat(self, food: str) -> None:
        pass

    @property
    @abstractmethod
    def can_walk(self) -> bool:
        pass


class Cat(Animal):
    def eat(self, food: str) -> None:
        ...  # Body omitted

    @property
    def can_walk(self) -> bool:
        return True


x = Animal()  # Error: 'Animal' is abstract due to 'eat' and 'can_walk'
y = Cat()  # OK

