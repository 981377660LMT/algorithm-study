# 哈夫曼编码(huffmanCode)
# 哈夫曼编码是一种变长编码，即不同的字符对应的编码长度不同，且没有一个编码是另一个编码的前缀。


import heapq
from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Mapping, Optional, Tuple, TypeVar


T = TypeVar("T")


class HuffmanTree(Generic[T]):
    def build(self, freq: Mapping[T, int]) -> Tuple[Optional["Node[T]"], DefaultDict[T, str]]:
        """根据频率表构建哈夫曼树、编码映射表."""
        root = self._build_huffman_tree(freq)
        codes = defaultdict(str)
        if root is not None:
            self._build_codes(root, "", codes)
        return root, codes

    def encode(self, text: Iterable[T], codes: Mapping[T, str]) -> str:
        """返回编码后的文本."""
        return "".join(codes[char] for char in text)

    def decode(self, encoded_text: str, root: Optional["Node[T]"]) -> List[T]:
        """返回解码后的文本."""
        decoded_chars = []
        node = root
        for bit in encoded_text:
            node = node.left if bit == "0" else node.right  # type: ignore
            if node.char is not None:  # type: ignore
                decoded_chars.append(node.char)  # type: ignore
                node = root
        return decoded_chars

    def _build_huffman_tree(self, freq: Mapping[T, int]) -> Optional["Node[T]"]:
        heap = [Node(freq, char) for char, freq in freq.items()]
        heapq.heapify(heap)
        while len(heap) > 1:
            node1 = heapq.heappop(heap)
            node2 = heapq.heappop(heap)
            merged = Node(node1.freq + node2.freq, left=node1, right=node2)
            heapq.heappush(heap, merged)
        return heap[0] if heap else None

    def _build_codes(
        self, node: Optional["Node[T]"], current_code: str, codes: DefaultDict[T, str]
    ) -> None:
        if node is None:
            return
        if node.char is not None:  # leaf node
            codes[node.char] = current_code
            return
        self._build_codes(node.left, current_code + "0", codes)
        self._build_codes(node.right, current_code + "1", codes)


class Node(Generic[T]):
    __slots__ = ("freq", "char", "left", "right")

    def __init__(
        self,
        freq: int,
        char: Optional[T] = None,  # !None for internal nodes, character for leaf nodes
        left: Optional["Node[T]"] = None,
        right: Optional["Node[T]"] = None,
    ):
        self.freq = freq
        self.char = char
        self.left = left
        self.right = right

    def __lt__(self, other: "Node[T]") -> bool:
        return self.freq < other.freq


if __name__ == "__main__":
    freq = {"a": 40, "m": 10, "l": 7, "f": 8, "t": 15}
    huffman = HuffmanTree()
    root, codes = huffman.build(freq)
    print(codes)
    encoded = huffman.encode("a" * 40 + "m" * 10 + "l" * 7 + "f" * 8 + "t" * 15, codes)
    print("Encoded Text:", encoded)
    decoded = huffman.decode(encoded, root)
    print("Decoded Text:", decoded)
