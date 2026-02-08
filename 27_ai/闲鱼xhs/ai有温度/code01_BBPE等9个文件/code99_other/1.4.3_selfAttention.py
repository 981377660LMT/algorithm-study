import torch
import torch.nn as nn
import torch.nn.functional as F

class SelfAttention(nn.Module):
    def __init__(self, embed_size): 
        # 输入的维度
        super(SelfAttention, self).__init__()
        self.embed_size = embed_size
        # QKV线性变换矩阵
        self.query = nn.Linear(embed_size, embed_size)
        self.Key = nn.Linear(embed_size, embed_size)
        self.value = nn.Linear(embed_size, embed_size)
        # 注意力权重dropout
        self.attention_dropout = nn.Dropout(0.1)
        # 输出线性变换
        self.output_linear = nn.Linear(embed_size, embed_size)

    def forward(self, datas, mask): # batch_size, seq_len, embed_size
        # QKV线性变换
        Q = self.query(datas)
        K = self.Key(datas)
        V = self.value(datas)
        # 矩阵乘法
        # Q :batch_size, seq_len, embed_size
        # K :batch_size, embed_size, seq_len
        P = torch.matmul(Q, K.transpose(-1, -2))
        # P: batch_size, seq_len, seq_len
        # P/sqrt(embed_size)，根号d，d是embed_size
        attention_scores = P/torch.sqrt(torch.tensor(self.embed_size))
        # 掩码
        if mask is not None:
            attention_scores = attention_scores.masked_fill(mask==0, float('-inf'))
        # 归一化
        attention_weight = F.softmax(attention_scores, dim = -1)
        # 注意力权重dropout
        attention_weight = self.attention_dropout(attention_weight)
        # 加权求和
        # attention_weight: batch_size, seq_len, seq_len
        # V: batch_size, seq_len, embed_size
        output = torch.matmul(attention_weight, V)
        # 输出线性变换
        output = self.output_linear(output)
        # output: batch_size, seq_len, embed_size
        return output

    
datas = torch.rand(3, 2, 4)
mask = torch.ones(3, 2, 2)
net = SelfAttention(4)
output = net(datas, mask)
print(output.shape)