from transformers import BertForSequenceClassification, BertTokenizer
from config import Config
import torch
import torch.nn as nn


# 奖励模型
class RewardModel(nn.Module):
    def __init__(self, config: Config):
        super(RewardModel, self).__init__()
        self.config = config
        self.reward_tokenizer = BertTokenizer.from_pretrained(self.config.Sentiment_model)
        self.reward_model = BertForSequenceClassification.from_pretrained(self.config.Sentiment_model).to(
            self.config.device)

    @torch.no_grad()
    def forward(self, text):
        input_ids, attention_mask = self.data_process(text)
        output = self.reward_model(torch.tensor(input_ids).to(self.config.device),
                                   torch.tensor(attention_mask).to(self.config.device))
        probs = torch.softmax(output.logits.clone().detach(), dim=1).tolist()
        reward = [prob[0] for prob in probs]
        return reward

    def data_process(self, texts):
        attention_mask = []
        input_ids = [self.reward_tokenizer.encode(text)[:512] for text in texts]
        max_length = max(len(i) for i in input_ids)
        res = []
        for one in input_ids:
            padding_num = max_length - len(one)
            res.append(one + [self.reward_tokenizer.pad_token_id] * padding_num)
            attention_mask.append([1] * len(one) + [0] * padding_num)
        return res, attention_mask

