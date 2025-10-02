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
