class ClockPageReplacement:
    def __init__(self, capacity: int):
        self.capacity = capacity
        self.pages = []
        self.use_bit = []
        self.pointer = 0

    def access_page(self, page: int):
        if page in self.pages:
            index = self.pages.index(page)
            self.use_bit[index] = 1
        else:
            if len(self.pages) < self.capacity:
                self.pages.append(page)
                self.use_bit.append(1)
            else:
                while self.use_bit[self.pointer] == 1:
                    self.use_bit[self.pointer] = 0
                    self.pointer = (self.pointer + 1) % self.capacity
                self.pages[self.pointer] = page
                self.use_bit[self.pointer] = 1
                self.pointer = (self.pointer + 1) % self.capacity

    def get_cache(self):
        return self.pages


# 示例使用
clock_cache = ClockPageReplacement(3)
pages = [1, 2, 3, 1, 4, 5]

for page in pages:
    clock_cache.access_page(page)
    print(f"访问页面 {page} 后的缓存状态: {clock_cache.get_cache()}")
