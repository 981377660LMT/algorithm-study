from typing import List


class A:
    def __init__(self, x: int) -> None:
        self.x = x  # Aha, attribute 'x' of type 'int'


a = A(1)
a.x = 2  # OK!
a.y = 3  # Error: "A" has no attribute "y"


# 可以使用类型注释在类体中显式声明`实例或者类变量`:
class B:
    ids: List[int]


b = B()
# b.ids = [1]  # OK
print(B.ids)

