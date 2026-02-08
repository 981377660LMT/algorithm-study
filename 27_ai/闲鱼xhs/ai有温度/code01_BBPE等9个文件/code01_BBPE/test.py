from bbpe import BBPETokenizer

# 这里可以换成llama3和qwen的词表进行对比
vocab_path = "./bbpe_models/斗破苍穹/vocab.json"
merges_path = "./bbpe_models/斗破苍穹/merges.txt"

tokenizer = BBPETokenizer(vocab_path, merges_path)

encode_text = input("请输入要编码的字符串：")
print(f"{'#' * 20}编码'{encode_text}'字符串的结果如下{'#' * 20}")

text_id = tokenizer.encode(encode_text)
id_text = tokenizer.decode(tokenizer.encode(encode_text))
text_token = tokenizer.tokenize(encode_text)

print(f"'{encode_text}'字符串编码之后对应的id: {text_id}")
print(f"'{encode_text}'字符串编码之后在反编码对应的字符串: {encode_text}")
print(f"'{encode_text}'字符串编码之后对应的token: {text_token}")

print('#' * 60)

# 批量编码
print("格式：list[str,……，str]，案例：['你好多谢支持', '请关注小红书@AI有温度']")
batch_encode_text = input("请输入要编码的字符串（按上面格式）：")
batch_encode_text = eval(batch_encode_text)
batch_encode_res = tokenizer.encode_batch(batch_encode_text, num_threads=2)

print(f"---------------------------编码批量数据结果如下---------------------------")
print(f"批量编码的数据'{batch_encode_text}'")
print(f"批量编码的结果如下：{batch_encode_res}")
