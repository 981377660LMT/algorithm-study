from torch import nn
import torch
import math


class EmbeddingWithPosition(nn.Module):
    def __init__(self, vocab_size, dim, seq_max_len):
        super().__init__()

        self.embedding = nn.Embedding(vocab_size, dim)
        position_idx = torch.arange(0, seq_max_len, dtype=torch.float).unsqueeze(-1)
        position_emb_fill = position_idx * torch.exp(
            -torch.arange(0, dim, 2) * math.log(10000.0) / dim
        )
        pos_encoding = torch.zeros(seq_max_len, dim)
        pos_encoding[:, 0::2] = torch.sin(position_emb_fill)
        pos_encoding[:, 1::2] = torch.cos(position_emb_fill)
        self.register_buffer("pos_encoding", pos_encoding)

    def forward(self, x):
        # x: (batch_size, seq_len)
        x = self.embedding(x)
        # x: (batch_size, seq_len, hidden_dim)
        x = (
            x + self.pos_encoding.unsqueeze(0)[:, : x.size()[1], :]
        )
        return x
