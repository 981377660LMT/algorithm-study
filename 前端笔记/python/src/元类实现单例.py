class SingletonMeta(type):
    """自定义单例元类"""

    def __init__(self, *args, **kwargs):
        self.__instance = None
        super().__init__(*args, **kwargs)

    def __call__(self, *args, **kwargs):
        if self.__instance is None:
            self.__instance = super().__call__(*args, **kwargs)
        return self.__instance


class President(metaclass=SingletonMeta):
    def __init__(self):
        print('Creating')


a = President()
b = President()
c = President()

