import torch
import torch.nn as nn
import torch.nn.functional as F

class GroupQueryAttention(nn.Module):
    # embed_size: 隐藏层维度
    # nums_head: 注意力头数
    # nums_kv_head: 计算注意力的 key value 头数
    
    # nums_kv_head == 1时，就是MQA
    def __init__(self, embed_size, nums_head, nums_kv_head):
        super(GroupQueryAttention, self).__init__()
        assert embed_size % nums_head == 0 # 可以整除
        assert nums_head % nums_kv_head == 0  # N 个 query head 为一组
        self.embed_size = embed_size
        self.nums_head = nums_head
        self.nums_kv_head = nums_kv_head
        self.head_dim = embed_size // nums_head

        # 初始化 qkv
        self.query = nn.Linear(embed_size, embed_size)  
        # k v out shape (nums_kv_head * head_dim)
        self.key = nn.Linear(embed_size, nums_kv_head * self.head_dim)
        self.value = nn.Linear(embed_size, nums_kv_head * self.head_dim)
        self.attention_dropout = nn.Dropout(0.1)
        self.output_linear = nn.Linear(embed_size, embed_size) 

    def forward(self, datas, mask=None):
        # X shape (batch, seq, hidden_dim)
        batch_size, seq, _ = datas.size()

        q = self.query(datas)  # （batch, seq, hidden_dim)
        k = self.key(datas)
        v = self.value(datas) 

        # attention_weight 目标shape 是 (batch, nums_head, seq, seq)
        q = q.view(batch_size, seq, self.nums_head, self.head_dim)
        k = k.view(batch_size, seq, self.nums_kv_head, self.head_dim)
        v = v.view(batch_size, seq, self.nums_kv_head, self.head_dim)

        # 关注: nums_head 和 nums_key_value_head 的关系
        q = q.transpose(1, 2) # (b, nums_head, seq, head_dim)
        k = k.transpose(1, 2) # (b, nums_kv_head, seq, head_dim)
        v = v.transpose(1, 2) # (b, nums_kv_head, seq, head_dim)

        # k v repeat 将自身复制并附加到自己的后面，在自身之后插入一个自身的副本
        # （batch, nums_key_value_head, seq, head_dim） -> （batch, nums_head, seq, head_dim）
        # tensor =torch.tensor([1,2,3])
        # tensor.repeat_interleave(2, dim=0)
        # [1,1,2,2,3,3] 
        k = k.repeat_interleave(self.nums_head // self.nums_kv_head, dim=1)
        v = v.repeat_interleave(self.nums_head // self.nums_kv_head, dim=1)

        P = torch.matmul(q, k.transpose(-2, -1))
        attention_scores = P/torch.sqrt(torch.tensor(self.head_dim))

        if mask is not None:
            # 将是0的，被掩盖的，填写成负无穷，在计算softmax时不考虑权重
            attention_scores = attention_scores.masked_fill(mask == 0, float('-inf'))

        attention_weight = F.softmax(attention_scores, dim=-1)
        attention_weight = self.attention_dropout(attention_weight)
        output = torch.matmul(attention_weight, v)  
        # (b, nums_head, seq, head_dim)
        # output(b, seq, hidden_dim)
        # 返回的张量在内存中是连续存储的
        output = output.transpose(1, 2).reshape(batch_size, seq, -1)
        output = self.output_linear(output)

        return output
# 测试
x = torch.rand(3, 2, 128)
# 8个头，4个key value 头， 显存为以前的一半
# mask: （batch_size, 1, seq_length, seq_length）
mask = torch.ones(3, 1, 2, 2)
net = GroupQueryAttention(embed_size = 128, nums_head = 8, nums_kv_head = 4)
print(net(x, mask).shape)
