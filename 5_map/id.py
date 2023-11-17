from collections import defaultdict


id_ = defaultdict(lambda: len(id_))


for i in range(10):
    print(id_[i])
