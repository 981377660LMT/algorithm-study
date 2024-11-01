# TODO

import time


class DataItem:
    def __init__(self, value):
        self.versions = [(value, time.time())]


class Transaction:
    def __init__(self, transaction_id):
        self.id = transaction_id
        self.start_time = time.time()
        self.read_set = set()
        self.write_set = {}


global_version = 0


def begin_transaction():
    transaction_id = int(time.time() * 1000)
    return Transaction(transaction_id)


def read(transaction, data_item):
    for value, version in reversed(data_item.versions):
        if version <= transaction.start_time:
            transaction.read_set.add(data_item)
            return value
    return None


def write(transaction, data_item, value):
    transaction.write_set[data_item] = value


def commit(transaction):
    global global_version
    for data_item in transaction.read_set:
        if data_item.versions[-1][1] > transaction.start_time:
            return False

    for data_item, value in transaction.write_set.items():
        global_version += 1
        data_item.versions.append((value, time.time()))

    return True


def rollback(transaction):
    transaction.read_set.clear()
    transaction.write_set.clear()


# 示例用法
data_item = DataItem(10)
transaction = begin_transaction()
print(read(transaction, data_item))  # 输出: 10
write(transaction, data_item, 20)
commit(transaction)
print(read(transaction, data_item))  # 输出: 20
