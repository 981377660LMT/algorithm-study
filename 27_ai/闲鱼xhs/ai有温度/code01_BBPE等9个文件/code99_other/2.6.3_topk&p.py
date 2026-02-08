def top_k_sampling(logits, k=50):
    # 过滤top-k
    topk_logits, topk_idx = logits.topk(k, dim=-1)
    
    # 采样
    probs = F.softmax(topk_logits, dim=-1)
    next_token_idx = torch.multinomial(probs, 1)
    return topk_idx.gather(-1, next_token_idx)

def top_p_sampling(logits, p=0.9):
    # 核采样
    sorted_logits, sorted_idx = torch.sort(logits, descending=True)
    cumulative_probs = torch.cumsum(F.softmax(sorted_logits, dim=-1), dim=-1)
    
    # 移除累计概率>p的token
    mask = cumulative_probs <= p
    mask[..., 0] = True  # 确保至少一个token
    
    filtered_logits = torch.where(mask, sorted_logits, torch.full_like(sorted_logits, -float('inf')))
    return torch.multinomial(F.softmax(filtered_logits, dim=-1), 1)