# 顺序违反
# !use_resource 线程可能在 initialize 线程初始化资源之前执行，导致 resource 未被初始化而出现错误。

import threading

resource = None
lock = threading.Lock()


def initialize():
    global resource
    with lock:
        resource = "资源已初始化"
        print("资源初始化完成")


def use_resource():
    with lock:
        if resource is not None:
            print(f"使用{resource}")
        else:
            print("资源未初始化")


t1 = threading.Thread(target=use_resource)
t2 = threading.Thread(target=initialize)

t1.start()
t2.start()

t1.join()
t2.join()

############################################
import threading

resource = None
resource_initialized = threading.Event()


def initialize():
    global resource
    resource = "资源已初始化"
    print("资源初始化完成")
    resource_initialized.set()


def use_resource():
    resource_initialized.wait()  # 等待初始化完成(await)
    print(f"使用{resource}")


t1 = threading.Thread(target=use_resource)
t2 = threading.Thread(target=initialize)

t1.start()
t2.start()

t1.join()
t2.join()
