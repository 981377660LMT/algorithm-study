// import operator

// class LazyObject:
//     @staticmethod
//     def new_method_proxy(func):
//         def inner(self, *args):
//             if not self._is_init:
//                 self._setup()
//             return func(self._wrapped, *args)

//         return inner

//     _resolved = None
//     _is_init = False

//     __getattr__ = new_method_proxy(getattr)
//     __bytes__ = new_method_proxy(bytes)
//     __str__ = new_method_proxy(str)
//     __bool__ = new_method_proxy(bool)
//     __dir__ = new_method_proxy(dir)
//     __hash__ = new_method_proxy(hash)
//     __class__ = property(new_method_proxy(operator.attrgetter("__class__")))  # type: ignore
//     __eq__ = new_method_proxy(operator.eq)
//     __lt__ = new_method_proxy(operator.lt)
//     __gt__ = new_method_proxy(operator.gt)
//     __ne__ = new_method_proxy(operator.ne)
//     __hash__ = new_method_proxy(hash)
//     __getitem__ = new_method_proxy(operator.getitem)
//     __setitem__ = new_method_proxy(operator.setitem)
//     __delitem__ = new_method_proxy(operator.delitem)
//     __iter__ = new_method_proxy(iter)
//     __len__ = new_method_proxy(len)
//     __contains__ = new_method_proxy(operator.contains)

//     def __init__(self, factory):
//         # Assign using __dict__ to avoid the setattr method.
//         self.__dict__["_factory"] = factory

//     def __setattr__(self, name, value):
//         # These are special names that are on the LazyObject.
//         # every other attribute should be on the wrapped object.
//         if name in {"_is_init", "_wrapped"}:
//             self.__dict__[name] = value
//         else:
//             if not self._is_init:
//                 self._setup()
//             setattr(self._resolved, name, value)

//     def __delattr__(self, name):
//         if name == "_wrapped":
//             raise TypeError("can't delete _wrapped.")
//         if not self._is_init:
//             self._setup()
//         delattr(self._resolved, name)

//     def _setup(self):
//         self._resolved = self._factory()
//         self._is_init = True
