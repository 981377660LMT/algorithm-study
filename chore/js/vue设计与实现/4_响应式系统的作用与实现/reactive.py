# TODO:
# 仅支持字典类型的响应式，没有实现数组方法的响应式
# 缺少批量更新和调度机制
# 没有处理循环依赖和无限递归问题


from typing import Callable, Optional


class Dep:
    """依赖收集器"""

    current = Optional[Callable[[], None]]  # 当前effect

    def __init__(self):
        self.subscribers = set()  # 依赖此数据的effect集合

    def depend(self):
        """收集依赖"""
        if Dep.current:
            self.subscribers.add(Dep.current)

    def notify(self):
        """通知更新"""
        # 创建副本避免在迭代过程中集合发生变化
        for effect in list(self.subscribers):
            effect()


class ReactiveDict(dict):
    """响应式字典对象"""

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._deps = {}  # 存储每个key的依赖收集器

    def _get_dep(self, key):
        """获取指定key的依赖收集器"""
        if key not in self._deps:
            self._deps[key] = Dep()
        return self._deps[key]

    def __getitem__(self, key):
        """重写字典的获取操作，进行依赖收集"""
        dep = self._get_dep(key)
        dep.depend()  # 收集依赖

        value = super().__getitem__(key)
        # 如果值是字典类型，需要递归转换为响应式
        if isinstance(value, dict) and not isinstance(value, ReactiveDict):
            value = reactive(value)
            super().__setitem__(key, value)
        return value

    def __setitem__(self, key, value):
        """重写字典的设置操作，触发更新通知"""
        # 如果是新值是字典，转换为响应式对象
        if isinstance(value, dict) and not isinstance(value, ReactiveDict):
            value = reactive(value)

        # 更新值
        super().__setitem__(key, value)

        # 通知依赖更新
        if key in self._deps:
            self._deps[key].notify()


def reactive(obj):
    """将普通对象转换为响应式对象"""
    if isinstance(obj, dict) and not isinstance(obj, ReactiveDict):
        return ReactiveDict(obj)
    return obj


def effect(fn):
    """创建响应式effect"""

    def wrapped_effect():
        prev_effect = Dep.current
        Dep.current = wrapped_effect
        try:
            return fn()
        finally:
            Dep.current = prev_effect

    wrapped_effect()
    return wrapped_effect


def watch(source, callback):
    """监听数据变化"""
    old_value = [None]

    def update():
        new_value = source()
        if new_value != old_value[0]:
            callback(new_value, old_value[0])
            old_value[0] = new_value

    effect(update)


if __name__ == "__main__":
    # 创建响应式数据
    state = reactive({"count": 0, "user": {"name": "Alice", "age": 25}, "todos": []})

    # 使用effect跟踪变化并更新UI
    def update_ui():
        print(f"UI更新: 计数={state['count']}, 用户={state['user']['name']}")

    effect_instance = effect(update_ui)

    # 使用watch监听变化
    def watch_count():
        return state["count"]

    def on_count_change(new_val, old_val):
        print(f"计数从 {old_val} 变为 {new_val}")

    watch(watch_count, on_count_change)

    # 测试数据变化
    print("\n修改count:")
    state["count"] = 1  # 触发UI更新和watch回调

    print("\n修改嵌套属性:")
    state["user"]["name"] = "Bob"  # 触发UI更新和计算属性重新计算
