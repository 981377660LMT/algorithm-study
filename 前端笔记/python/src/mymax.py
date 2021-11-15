# 用Python代码实现Python内置函数max。
def my_max(*args, key=None, default=None):
    """
    获取可迭代对象中最大的元素或两个及以上实参中最大的元素
    :param args: 一个可迭代对象或多个元素
    :param key: 提取用于元素比较的特征值的函数，默认为None
    :param default: 如果可迭代对象为空则返回该默认值，如果没有给默认值则引发ValueError异常
    :return: 返回可迭代对象或多个元素中的最大元素
    """
    if len(args) == 1 and len(args[0]) == 0:
        if default:
            return default
        else:
            raise ValueError('max() arg is an empty sequence')
    items = args[0] if len(args) == 1 else args
    max_elem, max_value = items[0], items[0]
    if key:
        max_value = key(max_value)
    for item in items:
        value = item
        if key:
            value = key(item)
        if value > max_value:
            max_elem, max_value = item, value
    return max_elem
