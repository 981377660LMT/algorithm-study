from typing import cast


a = cast(str, [4])

print(a.capitalize())
assert isinstance(a, int)
# AttributeError: 'list' object has no attribute 'capitalize'
a
# cast相当于ts里的id函数
# cast(typ: Type[_T@cast], val: Any) -> _T@cast
