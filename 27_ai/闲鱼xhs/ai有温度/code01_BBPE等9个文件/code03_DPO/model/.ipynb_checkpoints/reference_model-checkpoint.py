from transformers import AutoModelForCausalLM
import torch.nn as nn
import torch


class ReferenceModel(nn.Module):
    def __init__(self, config):
        super(ReferenceModel, self).__init__()
        self.config = config
        self.reference_model = AutoModelForCausalLM.from_pretrained(self.config.gpt_model, torch_dtype="auto").to(
            self.config.device)

    @torch.no_grad()
    def forward(self, input_ids, attention_mask):
        logits = self.reference_model(input_ids=input_ids,
                                      attention_mask=attention_mask).logits
        return logits


