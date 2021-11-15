# type元类动态创建类
# type(name, bases, dict)，接收三个参数

# 第一个参数name是指要创建类的名称
# 第二个参数bases是指需要继承父类的元组
# 第三个参数dict是类的属性


# 我们知道可以使用type(name, bases, dict)
# 来创建类，如果当使用type元类无法满足我们的一些需求时，我们可以自定义一个元类并使用该元类去创建类吗？
# 答案是可以的，下面我们来看一下：


class MyMetaClass(type):
    def __init__(cls, name, bases, dict):
        super().__init__(name, bases, dict)

        cls.int_attrs = {}

        for k, v in dict.items():
            if type(v) is int:
                cls.int_attrs[k] = v


User = MyMetaClass(
    'User', (), {'name': 'tigeriaf', "age": 24, "level": 2, "introduction": "Python全菜工程师"}
)
print(User)  # <class '__main__.User'>
user = User()
print(user.name)  # tigeriaf
print(user.int_attrs)  # {'age': 24, 'level': 2}
