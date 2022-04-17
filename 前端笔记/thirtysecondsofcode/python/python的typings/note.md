[容器的抽象基类](https://docs.python.org/zh-cn/3/library/collections.abc.html)

1. **Sequence** 是 Python 的一种内置类型（built-in type），内置类型就是构建在 Python Interpreter 里面的类型，三种基本的 Sequence Type 是 **list**（表），**tuple**（定值表，或翻译为元组），**range**（范围）
   这些类型的共同点是集合中的元素是有序排列的
2. 可调对象（Callable）
3. 迭代器与生成器与可迭代对象
   ABC|继承自|抽象方法|mixin 方法
   -----|-----|-----|-----
   Iterable|无|`__iter__`|
   Iterator|Iterable|`__next__`|`__iter__`
   Generator|Iterator|`send`， `throw`|`close`，`__iter__`，`__next__`
   `
4. Sequence 与 Iterable 区别
   iterable： 至少定义了`__iter__()`方法的对象。
   sequence： 至少定义了`__len__()`或者`__getitem__()`方法的对象。
   iterator：至少定义`__iter__()`和`__next__()`方法的对象。
5. ts 的 never 等价于 py 的 NoReturn （函数无终点）
   ts 的 void 等价于 py 的 None (不返回值的函数的返回类型)
   ts 的 any 等价于 py 的 Any
   ts 的 unkonwn 等价于 py 的 object
   ts 的 类型别名 等价于 py 的 直接变量赋值
   ts 的 联合类型 等价于 py 的 Union[]或者`|`
   ts 的 可选类型 等价于 py 的 `Optional[]`或者`|None`
   ts 的 typeof 等价于 py 的 `type(基本类型)`或者`Type(类)`
   ts 的 字面量类型 等价于 py 的 `Literal[]`
   ts 的 as const 等价于 py 的 `:Final[]`
   ts 的 类型守卫 等价于 py 的 `->TypeGuard[]`
   ts 的 record 等价于 py 的 `Mapping[]`
   ts 的 非 readonly 等价于 py 的 `Mutable`
   ts 的 泛型参数 等价于 py 的 `TypeVar()`
   ts 的 S extends string 等价于 py 的 `AnyStr 这里(AnyStr = TypeVar('AnyStr', str, bytes)。)`
   ts 的 类型断言 as 等价于 py 的 `cast(类型，变量)或者 assert isinstance(变量，类型)(后者是runtime check)`
   ts 的 interface 等价于 py 的 `Protocol`

```Python
from typing import NoReturn

def black_hole() -> NoReturn:
    raise Exception("There is no going back ...")
```

5. 实现与抽象（duckType)
   注解参数时,最好使用**抽象容器类型**
   泛型具象容器|泛型抽象容器|
   -----|-----|
   list|Sequence/Iterable
   dict|Mapping
   set|AbstractSet
6. 函数泛型的显示

```Python

X = TypeVar('X')
Y = TypeVar('Y')


def lookup_name(mapping: Mapping[X, Y], key: X, default: Y) -> Y:
    try:
        return mapping[key]
    except KeyError:
        return default

调用时显示：`泛型参数@函数名`
(mapping: Mapping[X@lookup_name, Y@lookup_name], key: X@lookup_name, default: Y@lookup_name) -> Y@lookup_name
```

7. typing 模块里的`_alias`函数

```Python
# Various ABCs mimicking those in collections.abc.
def _alias(origin, params, inst=True):
    return _GenericAlias(origin, params, special=True, inst=inst)

# NOTE: Mapping is only covariant in the value type.
Mapping = _alias(collections.abc.Mapping, (KT, VT_co))
MutableMapping = _alias(collections.abc.MutableMapping, (KT, VT))
Iterable = _alias(collections.abc.Iterable, T_co)
Set = _alias(set, T, inst=False)

# Some unconstrained type variables.  These are used by the container types.
# (These are not for export.)
T = TypeVar('T')  # Any type.
KT = TypeVar('KT')  # Key type.
VT = TypeVar('VT')  # Value type.
T_co = TypeVar('T_co', covariant=True)  # Any type covariant containers.
V_co = TypeVar('V_co', covariant=True)  # Any type covariant containers.
VT_co = TypeVar('VT_co', covariant=True)  # Value type covariant containers.
T_contra = TypeVar('T_contra', contravariant=True)  # Ditto contravariant.
# Internal type variable used for Type[].
CT_co = TypeVar('CT_co', covariant=True, bound=type)
covariant:协变
contravariant:逆变
```

8. None 是只有一个值 None 的类型。None 也用作不返回值的函数的返回类型，即隐式返回 None 的函数。
9. strictNullChecks：
   默认情况下，null 和 undefined 都可以分配给 ts 中的所有类型。
   在严格空值检查模式下, **值 null 和 undefined 不再属与所有类型并且只能赋值给它们自己对应的类型和 any** (一个例外是 undefined 也可以被复制给 void). 所以, 虽然在普通类型检查模式 T 和 T | undefined 意义相同 (因为 undefined 被认为是任何 T 的子类型), 在严格类型检查模式下它们是不同的, 并且只有 T | undefined 允许 undefined 作为值. T 和 T | null 的关系也是如此.

10. python 中如何实现 **interface**(约束类的公有实例属性)
    实际产生：使用类或者抽象类即可
    静态检查：`Protocol`
