import torch
import torch.nn as nn
import regex as re
import json
from collections import Counter
from concurrent.futures import ThreadPoolExecutor


def bytes_to_unicode():
    """
    返回utf-8字节列表和到unicode字符串的映射。我们特别避免映射到bbpe代码所依赖的空白/控制字符。
    可逆的bbpe代码在unicode字符串上工作。这意味着如果您想避免UNKs，您需要在您的词汇表中使用大量的unicode字符。
    当你有一个10B的token数据集时，你最终需要大约5K才能获得良好的覆盖。这是你正常情况下的一个显著比例，
    比如说，32K的词汇量。为了避免这种情况，我们希望查找表介于utf-8字节和unicode字符串之间。
    """
    bs = (
            list(range(ord("!"), ord("~") + 1)) + list(range(ord("¡"), ord("¬") + 1)) + list(
        range(ord("®"), ord("ÿ") + 1))
    )
    cs = bs[:]
    n = 0
    for b in range(2 ** 8):
        if b not in bs:
            bs.append(b)
            cs.append(2 ** 8 + n)
            n += 1
    cs = [chr(n) for n in cs]
    return dict(zip(bs, cs))


class BBPETokenizer(nn.Module):

    def __init__(self, vocab_path: str, merges_path: str):
        super().__init__()
        with open(vocab_path, "r", encoding="utf-8") as f:  # 获得词表
            vocab = json.load(f)
        with open(merges_path, "r", encoding="utf-8") as f:  # 获得合并token规则词表
            merges = f.read()

        # 将合并存储为元组列表，删除最后一个空白行
        merges = [tuple(merge_str.split()) for merge_str in merges.split("\n")[:-1]]

        # token到BBPE解码索引映射
        self.encoder = vocab
        self.decoder = {v: k for k, v in self.encoder.items()}

        # 字节到unicode字符映射，256个字符
        self.byte_encoder = bytes_to_unicode()
        self.byte_decoder = {v: k for k, v in self.byte_encoder.items()}

        self.bbpe_ranks = dict(zip(merges, range(len(merges))))
        self.cache = {}

        # 预标记化拆分正则表达式模式
        self.pat = re.compile(r"""
                                 's|'t|'re|'ve|'m|'ll|'d|  # 常见的收缩
                                 \ ?\p{L}+|\ ?\p{N}+|  # 可选空格，后跟1+ unicode字母或数字
                                 \ ?[^\s\p{L}\p{N}]+|  # 可选空格，后面跟着1+非空白/字母/数字
                                 \s+(?!\S)|  # 1+空白字符，后面没有非空白字符
                                 \s+  # 1+空格字符
                                 """, re.X)

    def forward(self, text):
        if isinstance(text, list):
            # 批量编码
            tokens = self.encode_batch(text)
            tokens = [token for row in tokens for token in row]
        else:
            # 编码字符串
            tokens = self.encode(text)
        return torch.tensor(tokens)

    def bbpe(self, token):
        '''
        对token应用合并规则
        '''
        if token in self.cache:
            return self.cache[token]

        chars = [i for i in token]
        # 对于每个合并规则，尝试合并任何相邻的字符对
        for pair in self.bbpe_ranks.keys():
            i = 0
            while i < len(chars) - 1:
                if chars[i] == pair[0] and chars[i + 1] == pair[1]:
                    chars = chars[:i] + ["".join(pair)] + chars[i + 2:]
                else:
                    i += 1
        self.cache[token] = chars
        return chars

    def encode(self, text: str) -> list[int]:
        '''
        将字符串编码为BBPE token
        '''
        bbpe_tokens_id = []
        # pattern使用要输入BBPE算法的正则表达式模式拆分文本
        for token in re.findall(self.pat, text):
            # 将token转换为其字节表示，将字节映射到其unicode表示
            token = "".join(self.byte_encoder[b] for b in token.encode("utf-8"))
            # 对token执行bbpe合并，然后根据编码器将结果映射到它们的bbpe索引
            bbpe_tokens_id.extend(self.encoder[bpe_token] for bpe_token in self.bbpe(token))
        return bbpe_tokens_id

    def tokenize(self, text):
        """
        获得编码后的字符
        :param text: 文本
        :return: 返回编码后的字符
        """
        bbpe_tokens = []
        # pattern使用要输入BBPE算法的正则表达式模式拆分文本
        for token in re.findall(self.pat, text):
            # 将token转换为其字节表示，将字节映射到其unicode表示
            token = "".join(self.byte_encoder[b] for b in token.encode("utf-8"))
            # 对token执行bbpe合并，然后根据编码器获得结果
            bbpe_tokens.extend(bpe_token for bpe_token in self.bbpe(token))
        return bbpe_tokens

    def encode_batch(self, batch: list[str], num_threads=4):
        '''
        将字符串列表编码为BBPE token列表
        '''
        with ThreadPoolExecutor(max_workers=num_threads) as executor:
            result = executor.map(self.encode, batch)
        return list(result)

    def decode(self, tokens) -> str:
        if isinstance(tokens, torch.Tensor):
            tokens = tokens.tolist()
        text = "".join([self.decoder[token] for token in tokens])
        text = bytearray([self.byte_decoder[c] for c in text]).decode("utf-8", errors="replace")
        return text

    @staticmethod
    def train_tokenizer(data, vocab_size, vocab_outfile=None, merges_outfile=None):
        """
        :param data: 训练文本
        :param vocab_size: 保留词表的大小
        :param vocab_outfile: 保存词表的文件名
        :param merges_outfile: 保存合并字节的词表
        """

        if vocab_size < 256:
            raise ValueError("vocab_size must be greater than 256")

        # 预标记数据
        byte_encoder = bytes_to_unicode()
        pat_str = r"'s|'t|'re|'ve|'m|'ll|'d| ?[\p{L}]+| ?[\p{N}]+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+"
        split_words = [
            [byte_encoder[b] for b in token.encode("utf-8")] for token in re.findall(pat_str, data)
        ]
        # 向词汇表中添加基本词汇
        vocab = set(byte_encoder.values())
        merges = []

        # 构建词汇表，直到满足所需的词汇量
        while len(vocab) < vocab_size:
            print(len(vocab))
            pair_freq = Counter()
            # 找出最常见的一对
            for split_word in split_words:
                pair_freq.update(zip(split_word[:-1], split_word[1:]))
            most_common_pair = pair_freq.most_common(1)[0][0]

            #  更新词汇表和合并列表
            new_token = most_common_pair[0] + most_common_pair[1]
            vocab.add(new_token)
            merges.append(most_common_pair)

            # 对数据执行合并
            new_split_words = []
            for split_word in split_words:
                i = 0
                new_word = []
                # 对于单词中的每个重字符，尝试合并
                while i < len(split_word) - 1:
                    if (split_word[i], split_word[i + 1]) == most_common_pair:
                        new_word.append(new_token)
                        i += 2
                    else:
                        new_word.append(split_word[i])
                        i += 1
                if i == len(split_word) - 1:
                    new_word.append(split_word[i])
                new_split_words.append(new_word)
            split_words = new_split_words

        vocab = sorted(list(vocab))
        # 保存文件
        if merges_outfile != None:
            with open(merges_outfile, "w", encoding="utf-8") as f:
                for merge in merges:
                    f.write(merge[0] + " " + merge[1] + "\n")
        if vocab_outfile != None:
            with open(vocab_outfile, "w", encoding="utf-8") as f:
                json.dump({v: i for i, v in enumerate(vocab)}, f, ensure_ascii=False)



