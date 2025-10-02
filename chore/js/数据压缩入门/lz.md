好的，我们来详细讲解 LZ77 算法，并提供一个清晰的 Python 代码实现。LZ77 是整个 Lempel-Ziv 算法家族的基石，也是 `Gzip`, `ZIP`, `PNG` 等现代压缩格式核心思想的源头。

### LZ77 算法核心思想

LZ77 的天才之处在于它使用**滑动窗口 (Sliding Window)** 来动态地创建字典，从而避免了需要额外存储一个静态字典的开销。

1.  **滑动窗口**: 想象一个固定大小的窗口在你的数据上从左到右滑动。这个窗口被当前处理的位置（光标）分为两部分：

    - **搜索缓冲区 (Search Buffer)**: 位于光标**之前**的数据。这是我们的“字典”，包含了最近处理过的所有内容。我们在这里寻找重复。
    - **前向缓冲区 (Look-ahead Buffer)**: 位于光标**处及之后**的数据。这是我们希望进行压缩的“待匹配”内容。

2.  **匹配与编码**: 算法的核心循环如下：

    - 从前向缓冲区中取一段字符串。
    - 在搜索缓冲区中，从后向前（即从离光标最近的地方开始）查找这段字符串的最长匹配。
    - **如果找到匹配**:
      - 不输出原始字符串，而是输出一个**引用指针**，通常表示为 `(距离, 长度)`。
      - **距离 (Distance)**: 从当前光标位置回退多少个字符可以找到匹配串的开头。
      - **长度 (Length)**: 匹配的字符串有多长。
      - 然后，将光标向前移动 `长度` 个位置。
    - **如果找不到匹配**:
      - 无法压缩，只能输出一个**字面值 (Literal)**，即当前光标处的单个字符。
      - 然后，将光标向前移动 `1` 个位置。

3.  **输出**: 最终的压缩数据是一个由 `(距离, 长度)` 对和 `字面值` 混合组成的序列。这个序列之后可以再被霍夫曼编码等熵编码器进一步压缩，以达到更高的压缩率。

---

### Python 代码实现

下面的 Python 代码清晰地实现了 LZ77 的压缩和解压逻辑。为了教学目的，我们定义了 `WINDOW_SIZE` (搜索缓冲区大小) 和 `LOOKAHEAD_BUFFER_SIZE` (在前向缓冲区中尝试匹配的最大长度)。

```python
def find_longest_match(data, cursor, window_size, lookahead_buffer_size):
    """
    在搜索缓冲区中为前向缓冲区寻找最长的匹配。

    返回: (距离, 长度) 元组。如果未找到匹配，则返回 (0, 0)。
    """
    # 确定搜索缓冲区的实际范围
    search_buffer_start = max(0, cursor - window_size)
    search_buffer = data[search_buffer_start:cursor]

    best_match_length = 0
    best_match_distance = 0

    # 遍历前向缓冲区中所有可能的匹配长度，从最长到最短
    # 确保不超出数据末尾
    max_possible_length = min(len(data) - cursor, lookahead_buffer_size)
    for length in range(max_possible_length, 0, -1):
        # 获取要匹配的目标字符串
        target = data[cursor : cursor + length]

        # 在搜索缓冲区中从后向前查找
        # str.rfind() 是一个高效的查找方式
        # 返回值是子字符串在搜索缓冲区中的索引
        pos = search_buffer.rfind(target)

        if pos != -1:
            # 如果找到匹配
            distance = len(search_buffer) - pos
            best_match_length = length
            best_match_distance = distance
            # 因为我们是从最长长度开始找的，所以第一个找到的就是最长的
            return (best_match_distance, best_match_length)

    # 如果循环结束都没有找到任何匹配
    return (0, 0)

def compress(data, window_size=4096, lookahead_buffer_size=15):
    """
    使用 LZ77 算法压缩数据。

    返回: 一个由 (距离, 长度) 元组和字面值字符组成的列表。
    """
    compressed_data = []
    cursor = 0
    data_len = len(data)

    while cursor < data_len:
        distance, length = find_longest_match(data, cursor, window_size, lookahead_buffer_size)

        if length > 0:
            # 找到了一个匹配
            # 添加 (距离, 长度) 对到压缩结果中
            compressed_data.append((distance, length))
            # 将光标向前移动匹配的长度
            cursor += length
        else:
            # 没有找到匹配
            # 添加字面值到压缩结果中
            compressed_data.append(data[cursor])
            # 将光标向前移动 1
            cursor += 1

    return compressed_data

def decompress(compressed_data):
    """
    解压由 LZ77 算法压缩的数据。
    """
    decompressed_data = ""

    for token in compressed_data:
        if isinstance(token, tuple):
            # token 是一个 (距离, 长度) 对
            distance, length = token

            # 从已解压的数据末尾回溯 'distance' 个位置开始复制
            start_index = len(decompressed_data) - distance
            for i in range(length):
                # 复制字符。注意处理重叠复制的情况，
                # 例如 (1, 5) 用于 "abababa"
                decompressed_data += decompressed_data[start_index + i]
        else:
            # token 是一个字面值
            decompressed_data += token

    return decompressed_data

# --- 详细讲解与示例 ---
if __name__ == "__main__":
    # 一个具有很多重复的示例文本
    original_data = "abracadabra_abracadabra"
    print(f"原始数据: {original_data}")
    print(f"原始长度: {len(original_data)}\n")

    # 1. 压缩
    # 为了方便观察，我们使用一个较小的窗口
    compressed_output = compress(original_data, window_size=10, lookahead_buffer_size=5)
    print(f"压缩后的输出 (Token序列):")
    print(compressed_output)
    # 预期输出: ['a', 'b', 'r', 'a', 'c', 'a', 'd', (4, 4), '_', (12, 11)]
    # 让我们一步步分析这个输出是如何得到的：
    # 'a', 'b', 'r', 'a', 'c', 'a', 'd' -> 初始时找不到匹配，都是字面值。
    #   此时光标在第7个'a'处，搜索缓冲区是 "abracad"。
    #   前向缓冲区是 "abra..."。在 "abracad" 中找到了 "abra" 的匹配。
    #   距离是4 (从光标回退4个字符到'a')，长度是4。所以输出 (4, 4)。光标前进4。
    # '_' -> 找不到匹配，输出字面值。
    #   此时光标在第二个'a'处，搜索缓冲区是 "cadabra_ab"。
    #   前向缓冲区是 "abracadabra"。在搜索缓冲区和之前更早的数据中找到了 "abracadabra" 的匹配。
    #   距离是12 (从当前光标回退12个字符到第一个'a')，长度是11。所以输出 (12, 11)。

    print("\n--- 压缩过程逐步分析 ---")
    # 这是一个更简单的例子来展示逻辑
    simple_data = "aabcbbabc"
    print(f"分析简单数据: '{simple_data}'")
    # 预期: ['a', (1, 1), 'b', 'c', 'b', (4, 3)]
    # 1. 'a': 字面值。已处理: "a"
    # 2. 'a': 在"a"中找到匹配。距离1, 长度1 -> (1, 1)。已处理: "aa"
    # 3. 'b': 字面值。已处理: "aab"
    # 4. 'c': 字面值。已处理: "aabc"
    # 5. 'b': 字面值。已处理: "aabcb"
    # 6. 'abc': 在"aabcb"中找到"abc"的匹配。距离4, 长度3 -> (4, 3)。已处理: "aabcbabc"
    print(f"压缩结果: {compress(simple_data, 10, 5)}\n")


    # 2. 解压
    decompressed_result = decompress(compressed_output)
    print(f"解压后的数据: {decompressed_result}")
    print(f"解压后长度: {len(decompressed_result)}\n")

    # 3. 验证
    assert original_data == decompressed_result
    print("验证成功: 原始数据与解压后的数据完全一致！")

```

1.  **`find_longest_match` 函数**:

    - 这是算法的核心。它接收完整的 `data` 和当前 `cursor` 位置。
    - `search_buffer_start` 和 `search_buffer` 确定了我们的“字典”范围。
    - 循环从**最长**可能匹配的长度开始递减。这是为了确保我们找到的是最长的匹配，这是一种贪心策略，通常效果很好。
    - `data[cursor : cursor + length]` 是我们想要在前向缓冲区中匹配的目标。
    - `search_buffer.rfind(target)` 是一个非常关键的优化。`rfind` 从右到左查找，能更快地找到离光标最近的匹配，这通常能产生更小的“距离”值，有利于后续的熵编码。
    - 一旦找到第一个匹配（因为是从长到短搜索的，所以第一个就是最长的），立即计算 `distance` 和 `length` 并返回。
    - 如果循环结束都没有找到任何匹配，返回 `(0, 0)` 作为未找到的信号。

2.  **`compress` 函数**:

    - 这是主驱动函数。它使用一个 `while` 循环和 `cursor` 来遍历整个输入数据。
    - 在循环的每一步，它都调用 `find_longest_match`。
    - 如果返回的 `length > 0`，说明找到了匹配。它将 `(distance, length)` 元组存入结果，并把光标前进 `length`。
    - 否则，说明没有找到匹配。它将当前光标处的单个字符（字面值）存入结果，并把光标前进 `1`。
    - 最终返回一个由元组和字符混合组成的列表。

3.  **`decompress` 函数**:
    - 解压过程非常直接。它遍历压缩后的 `token` 列表。
    - 如果 `token` 是一个元组 `(distance, length)`，它就从已经解压出的字符串末尾，回退 `distance` 个位置，然后从那里开始复制 `length` 个字符。
      - **注意**: `decompressed_data[start_index + i]` 的实现可以正确处理重叠复制。例如，要解压 `(1, 5)` 并且已解压部分是 `"ab"`，它会先复制 `b` 变成 `"abb"`，然后复制第二个 `b` 变成 `"abbb"`，以此类推，最终正确生成 `"abbbbb"`。
    - 如果 `token` 是一个字符，就直接追加到结果字符串中。

这个实现清晰地展示了 LZ77 算法的“模型”部分——即将重复的字符串序列转换为 `(距离, 长度)` 对。在实际应用中，这个 `token` 序列会经过更复杂的处理，被转换成二进制流，并通常会用霍夫曼编码或 ANS 进行熵编码，以实现最终的文件体积压缩。
