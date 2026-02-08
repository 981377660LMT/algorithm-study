import torch
import torch.nn.functional as F


class Tools:
    def __init__(self, response_shape, response_ids_mask):
        """
        :param response_shape: 模型生成的句子长度
        :param response_ids_mask: 模型批量生成的句子中有padding，这里去除padding数据
        """
        self.response_shape = response_shape
        self.response_ids_mask = response_ids_mask

    def filter_mask(self, values):
        """
        :param values: 一般是prob_old、prob_ref、value(价值)的值
        :return: 去除padding之后的数据
        """
        return [value[-self.response_shape:][one_response_ids_mask] for value, one_response_ids_mask in
                zip(values, self.response_ids_mask)]

    @staticmethod
    def probs_from_logits(logits, labels):
        log_probs = F.log_softmax(logits, dim=2)
        probs = torch.gather(log_probs, 2, labels.unsqueeze(2)).squeeze(-1)
        return probs
