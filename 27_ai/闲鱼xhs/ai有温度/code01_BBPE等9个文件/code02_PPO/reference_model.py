from transformers import AutoModelForCausalLM
from tools import Tools
import torch.nn as nn
import torch


class ReferenceModel(nn.Module):
    def __init__(self, config):
        super(ReferenceModel, self).__init__()
        self.config = config
        self.reference_model = AutoModelForCausalLM.from_pretrained(self.config.gpt_model, torch_dtype="auto").to(
            self.config.device)

    @torch.no_grad()
    def forward(self, input_ids, attention_mask, tools: Tools):
        logits = self.reference_model(input_ids=input_ids,
                                      attention_mask=attention_mask).logits
        prob_refs = tools.probs_from_logits(logits[:, :-1, :], input_ids[:, 1:])
        prob_refs = tools.filter_mask(prob_refs)
        return prob_refs


