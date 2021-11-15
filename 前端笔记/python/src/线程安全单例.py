import threading
from functools import wraps


# synchronized关键字可以保证被它修饰的方法或者代码块在任意时刻只能有一个线程执行。
def synchronized(func):
    func.__lock__ = threading.Lock()

    @wraps(func)
    def wrapper(*args, **kwargs):
        with func.__lock__:
            return func(*args, **kwargs)

    return wrapper


class Singleton(object):
    _instance = None

    @synchronized
    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance

    def __init__(self, x, y):
        print('init successfully')
        self.x = x
        self.y = y


def func():
    obj = Singleton(1, 2)
    print(id(obj))


if __name__ == '__main__':
    for i in range(10):
        thread = threading.Thread(target=func)
        thread.start()

