import torch
from torch.utils.data import Dataset, DataLoader
import json

class CustomDataset(Dataset):
    def __init__(self, data_file, actor_tokenizer):
        self.actor_tokenizer = actor_tokenizer
        with open(data_file, 'r', encoding="utf-8") as f:
            self.data = json.load(f)
        self.total_samples = len(self.data)
        self.actor_padding_id = actor_tokenizer.pad_token_id

    def __len__(self):
        return self.total_samples

    def __getitem__(self, idx):
        # 根据索引加载数据
        # 这里可以根据需要从文件中读取数据并进行预处理
        line = self.data[idx]
        query = line["query"]
        system_content = line["system_content"]
        messages = [
            {"role": "system", "content": system_content},
            {"role": "user", "content": query},
        ]
        text = self.actor_tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=True
        )
        model_inputs = self.actor_tokenizer([text], return_tensors="pt")

        return [model_inputs.input_ids.tolist()[0], model_inputs.attention_mask.tolist()[0]]

    def collate_fn(self, batch):
        max_length = max([len(i[0]) for i in batch])
        input_ids = []
        mask_attention = []
        for one in batch:
            padding_num = max_length - len(one[0])
            input_ids.append([self.actor_padding_id] * padding_num + one[0])
            mask_attention.append([0] * padding_num + one[1])
        return torch.tensor(input_ids), torch.tensor(mask_attention)


if __name__ == '__main__':
    from config import Config
    from transformers import AutoTokenizer

    config = Config()
    tokenizer = AutoTokenizer.from_pretrained(config.gpt_model)
    # 创建自定义数据集实例
    dataset = CustomDataset(config.data_path, tokenizer)

    # 创建数据加载器并指定批次大小
    batch_size = 2
    data_loader = DataLoader(dataset, batch_size=batch_size, shuffle=True, collate_fn=dataset.collate_fn)

    # 使用生成器函数按需读取数据
    for batch in data_loader:
        print()
        # 在每个批次中进行模型训练
        # batch 包含了一个批次的样本数据
        # 在这里执行模型训练操作
        pass
