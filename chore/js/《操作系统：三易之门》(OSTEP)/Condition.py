import threading

# 共享资源
buffer = []
buffer_condition = threading.Condition()
BUFFER_SIZE = 5


def producer():
    while True:
        item = produce_item()
        with buffer_condition:
            while len(buffer) >= BUFFER_SIZE:
                buffer_condition.wait()
            buffer.append(item)
            buffer_condition.notify_all()


def consumer():
    while True:
        with buffer_condition:
            while not buffer:
                buffer_condition.wait()
            item = buffer.pop(0)
            buffer_condition.notify_all()
        consume_item(item)


def produce_item():
    # 生产项的逻辑
    return 1


def consume_item(item):
    # 消费项的逻辑
    pass


threading.Thread(target=producer).start()
threading.Thread(target=consumer).start()
