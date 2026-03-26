https://github.com/waylandzhang/Transformer-from-scratch
整体结构

```
Model (完整模型)
├── Token Embedding (词嵌入)
├── Positional Encoding (位置编码)
├── N × TransformerBlock (多个 Transformer 块)
│   ├── LayerNorm
│   ├── Multi-Head Attention
│   ├── LayerNorm
│   └── Feed Forward Network
├── Final LayerNorm (最后的归一化)
└── Output Linear (输出投影到词表)
```

- 为什么用 register_buffer？

Mask 不是参数（不需要训练），但需要跟模型一起移动到 GPU。register_buffer 就是专门干这个的。
