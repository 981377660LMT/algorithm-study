import torch
from torch.utils.data import Dataset, DataLoader
import json

class CustomDataset(Dataset):
    def __init__(self, data_file, tokenizer):
        self.tokenizer = tokenizer
        with open(data_file, 'r', encoding="utf-8") as f:
            self.data = json.load(f)
        self.total_samples = len(self.data)


    def __len__(self):
        return self.total_samples

    def __getitem__(self, idx):
        # 根据索引加载数据
        # 这里可以根据需要从文件中读取数据并进行预处理
        line = self.data[idx]
        query = line["input"]
        instruction = line["instruction"]
        rejected = line["output"]["rejected"]
        chosen = line["output"]["chosen"]
        messages = [
            {"role": "system", "content": "你是一个非常有帮助和智能的助手。"},
            {"role": "instrution", "content": instruction},
            {"role": "user", "content": query},
        ]
        text = self.tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=True
        )
        prompt_inputs = self.tokenizer.encode(text=text)
        rejected_inputs = self.tokenizer.encode(text=rejected)
        chosen_inputs = self.tokenizer.encode(text=chosen)
        return [prompt_inputs, rejected_inputs, chosen_inputs]

    def collate_fn(self, batch):
        def statisc_ids(data, labels_mask):
            max_length = max([len(i) for i in data])
            attention_masks = []
            return_ids = []
            labels_masks = []
            for one, mask in zip(data, labels_mask):
                padding_num = max_length - len(one)
                return_ids.append(one + [self.tokenizer.pad_token_id] * padding_num)
                labels_masks.append(mask + [0] * padding_num)
                attention_masks.append([1] * len(one) + [0] * padding_num)
            return return_ids, attention_masks, labels_masks

        inputs_ids = []
        labels_mask = []
        for i in range(1, 3):
            for one, prompt_input in zip(batch, [one[0] for one in batch]):
                res = one[i]
                inputs_ids.append(prompt_input + res)
                labels_mask.append([0] * len(prompt_input) + [1] * len(res))
        inputs_ids, inputs_masks, labels_masks = statisc_ids(inputs_ids, labels_mask)
        return {"inputs_ids": torch.tensor(inputs_ids), "inputs_masks": torch.tensor(inputs_masks),
                "labels_mask": torch.tensor(labels_masks)}


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
