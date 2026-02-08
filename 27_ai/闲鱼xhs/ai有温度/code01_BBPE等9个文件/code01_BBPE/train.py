from bbpe import BBPETokenizer

# 读取训练文件
with open("./train_data/斗破苍穹.txt", "r", encoding="utf-8") as f:
    data = f.read()
# 训练模型
vocab_size = 5000  # 词表大小
vocab_outfile = "vocab.json"  # 保存词表文件名
merges_outfile = "merges.txt"  # 保存合并字节的词表)
BBPETokenizer.train_tokenizer(data, vocab_size, vocab_outfile=vocab_outfile, merges_outfile=merges_outfile)
