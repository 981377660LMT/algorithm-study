from functools import wraps
from sys import hexversion

if hexversion < 0x03000000:
    print(1)
    from itertools import imap as map  # pylint: disable=redefined-builtin
    from itertools import izip as zip  # pylint: disable=redefined-builtin

    try:
        from thread import get_ident
    except ImportError:
        from dummy_thread import get_ident
else:
    from functools import reduce

    try:
        from _thread import get_ident
    except ImportError:
        from _dummy_thread import get_ident


# !修饰对象的 repr 方法，在递归调用时返回指定的填充值
def recursive_repr(fillvalue="..."):
    "Decorator to make a repr function return fillvalue for a recursive call."
    # pylint: disable=missing-docstring
    # Copied from reprlib in Python 3
    # https://hg.python.org/cpython/file/3.6/Lib/reprlib.py

    def decorating_function(user_function):
        repr_running = set()  # 环检测

        @wraps(user_function)  # 保留原始函数的元信息
        def wrapper(self):
            key = id(self), get_ident()  # 生成一个唯一的 key，由对象的 id 和当前线程的标识符（使用 get_ident() 函数）组成
            if key in repr_running:
                return fillvalue
            repr_running.add(key)
            try:
                result = user_function(self)
            finally:
                repr_running.discard(key)
            return result

        return wrapper

    return decorating_function
