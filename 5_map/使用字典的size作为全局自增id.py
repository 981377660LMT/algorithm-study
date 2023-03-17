# 使用字典的size作为全局自增id


pool = dict()
history = []
for name in ["Alice", "Bob", "Alice", "Sam", "Bob", "Green"]:
    pool.setdefault(name, len(pool))  # 全局自增id(len(pool)作为id)
    history.append(pool[name])


print(history)
