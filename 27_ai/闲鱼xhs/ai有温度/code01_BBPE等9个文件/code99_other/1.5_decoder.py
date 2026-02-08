import math
import torch
import torch.nn as nn

import warnings
warnings.filterwarnings(action="ignore")

# 写一个 Block
class SimpleDecoder(nn.Module):
    def __init__(self, hidden_dim, nums_head, dropout=0.1):
        super().__init__()

        self.nums_head = nums_head
        self.head_dim = hidden_dim // nums_head

        self.dropout = dropout

        # 这里按照 transformers 中的 decoder 来写，用 post_norm 的方式实现，注意有残差链接
        self.layernorm_att = nn.LayerNorm(hidden_dim, eps=0.00001)
        self.q_proj = nn.Linear(hidden_dim, hidden_dim)
        self.k_proj = nn.Linear(hidden_dim, hidden_dim)
        self.v_proj = nn.Linear(hidden_dim, hidden_dim)
        self.o_proj = nn.Linear(hidden_dim, hidden_dim)
        self.drop_att = nn.Dropout(self.dropout)

        # for ffn 准备
        self.up_proj = nn.Linear(hidden_dim, hidden_dim * 4)
        self.down_proj = nn.Linear(hidden_dim * 4, hidden_dim)
        self.layernorm_ffn = nn.LayerNorm(hidden_dim, eps=0.00001)
        self.act_fn = nn.ReLU()
        self.drop_ffn = nn.Dropout(self.dropout)

    def attention_output(self, query, key, value, attention_mask=None):
        # 计算两者相关性
        key = key.transpose(2, 3)  # (batch, num_head, head_dim, seq)
        att_weight = torch.matmul(query, key) / math.sqrt(self.head_dim)

        # attention mask 进行依次调整；变成 causal_attention
        if attention_mask is not None:
            # 变成下三角矩阵
            attention_mask = attention_mask.tril()
            att_weight = att_weight.masked_fill(attention_mask == 0, float("-1e20"))
        else:
            # 人工构造一个下三角的 attention mask
            attention_mask = torch.ones_like(att_weight).tril()
            att_weight = att_weight.masked_fill(attention_mask == 0, float("-1e20"))

        att_weight = torch.softmax(att_weight, dim=-1)
        print(att_weight)

        att_weight = self.drop_att(att_weight)

        mid_output = torch.matmul(att_weight, value)
        # mid_output shape is: (batch, nums_head, seq, head_dim)

        mid_output = mid_output.transpose(1, 2).contiguous()
        batch, seq, _, _ = mid_output.size()
        mid_output = mid_output.view(batch, seq, -1)
        output = self.o_proj(mid_output)
        return output

    def attention_block(self, X, attention_mask=None):
        batch, seq, _ = X.size()
        query = self.q_proj(X).view(batch, seq, self.nums_head, -1).transpose(1, 2)
        key = self.k_proj(X).view(batch, seq, self.nums_head, -1).transpose(1, 2)
        value = self.v_proj(X).view(batch, seq, self.nums_head, -1).transpose(1, 2)

        output = self.attention_output(
            query,
            key,
            value,
            attention_mask=attention_mask,
        )
        return self.layernorm_att(X + output)

    def ffn_block(self, X):
        up = self.act_fn(
            self.up_proj(X),
        )
        down = self.down_proj(up)
        # 执行 dropout
        down = self.drop_ffn(down)
        # 进行 norm 操作
        return self.layernorm_ffn(X + down)

    def forward(self, X, attention_mask=None):
        # X 一般假设是已经经过 embedding 的输入， (batch, seq, hidden_dim)
        # attention_mask 一般指的是 tokenizer 后返回的 mask 结果，表示哪些样本需要忽略
        # shape 一般是： (batch, nums_head, seq)

        att_output = self.attention_block(X, attention_mask=attention_mask)
        ffn_output = self.ffn_block(att_output)
        return ffn_output


# 测试

x = torch.rand(3, 4, 64)
net = SimpleDecoder(64, 8)
mask = (
    torch.tensor([[1, 1, 1, 1], [1, 1, 0, 0], [1, 1, 1, 0]])
    .unsqueeze(1)
    .unsqueeze(2)
    .repeat(1, 8, 4, 1)
)

net(x, mask).shape