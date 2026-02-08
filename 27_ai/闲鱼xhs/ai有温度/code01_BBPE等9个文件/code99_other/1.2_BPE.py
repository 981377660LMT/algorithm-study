from collections import defaultdict

def get_stats(vocab):
    """统计词汇表中相邻字符对的频率"""
    pairs = defaultdict(int)
    for word, freq in vocab.items():
        symbols = word.split()
        for i in range(len(symbols) - 1):
            pairs[(symbols[i], symbols[i + 1])] += freq
    return pairs

def merge_vocab(pair, vocab):
    """将词汇表中最频繁的字符对合并"""
    v_out = {}
    bigram = ' '.join(pair)
    replacement = ''.join(pair)
    for word in vocab:
        new_word = word.replace(bigram, replacement)
        v_out[new_word] = vocab[word]
    return v_out

def train_bpe(corpus, num_merges):
    """训练BPE模型，返回合并规则"""
    vocab = {}
    for word, freq in corpus.items():
        chars = list(word)
        vocab[' '.join(chars)] = freq
    
    merge_rules = []
    for i in range(num_merges):
        pairs = get_stats(vocab)
        if not pairs:
            break
        
        best_pair = max(pairs, key=pairs.get)
        merge_rules.append(best_pair)
        vocab = merge_vocab(best_pair, vocab)
    
    return merge_rules

def apply_bpe(word, merge_rules):
    """将训练好的BPE规则应用到单个单词"""
    if not word:
        return []
    
    # 将单词拆分为字符
    symbols = list(word)
    
    # 应用所有合并规则
    for pair in merge_rules:
        i = 0
        while i < len(symbols) - 1:
            if symbols[i] == pair[0] and symbols[i+1] == pair[1]:
                # 合并这两个符号
                symbols = symbols[:i] + [pair[0] + pair[1]] + symbols[i+2:]
            else:
                i += 1
    
    return symbols

# 示例用法
if __name__ == "__main__":
    # 用一些句子训练BPE模型
    # 实际训练中这个句子列表应该是从训练语料库中抽取的，就是一个个txt文件合并成一个训练语料库
    training_sentences = [
        "the quick brown fox",
        "the quick brown dog",
        "the slow blue cat",
        "the fast red bird",
        "a quick brown fox jumps"
    ]
    
    # 构建训练语料库（单词频率统计）
    corpus = defaultdict(int)
    for sentence in training_sentences:
        for word in sentence.split():
            corpus[word] += 1
    
    # 训练BPE模型，执行10次合并
    num_merges = 10
    merge_rules = train_bpe(corpus, num_merges)
    
    print(f"训练的合并规则 ({len(merge_rules)} 条):")
    for i, rule in enumerate(merge_rules, 1):
        print(f"{i}. {rule[0]} + {rule[1]} → {rule[0]}{rule[1]}")
    
    # 对新句子应用BPE分词
    test_sentence = "a fast brown fox jumps high"
    print(f"\n测试句子: {test_sentence}")
    print("BPE分词结果:")
    
    for word in test_sentence.split():
        bpe_tokens = apply_bpe(word, merge_rules)
        print(f"  {word}: {bpe_tokens}")