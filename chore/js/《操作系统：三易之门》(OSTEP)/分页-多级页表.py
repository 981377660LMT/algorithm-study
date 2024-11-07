class MultiLevelPageTable:
    def __init__(self, levels):
        self.levels = levels
        self.table = {}

    def map(self, virtual_address, physical_address):
        address = self._split_address(virtual_address)
        table = self._get_table(address, create=True)
        table[address[-1]] = physical_address

    def lookup(self, virtual_address):
        address = self._split_address(virtual_address)
        table = self._get_table(address)
        if table is None or address[-1] not in table:
            return None
        return table[address[-1]]

    def _get_table(self, address, create=False):
        table = self.table
        for level in range(self.levels - 1):
            if address[level] not in table:
                if create:
                    table[address[level]] = {}
                else:
                    return None
            table = table[address[level]]
        return table

    def _split_address(self, address):
        # 假设每级地址长度相同
        parts = []
        for _ in range(self.levels):
            parts.append(address & 0xFF)
            address >>= 8
        return parts[::-1]


# 示例使用
page_table = MultiLevelPageTable(levels=3)

# 映射虚拟地址到物理地址
page_table.map(0x123456, 0xABCDEF)


# 查找虚拟地址对应的物理地址
physical_address = page_table.lookup(0x123456)
if physical_address is not None:
    print(f"Physical Address: {hex(physical_address)}")
else:
    print("Address not mapped")
