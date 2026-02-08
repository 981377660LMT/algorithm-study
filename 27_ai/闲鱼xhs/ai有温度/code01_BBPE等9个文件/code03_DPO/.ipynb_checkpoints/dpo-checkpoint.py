import torch
import torch.nn.functional as F


class DPO:
    def __init__(self, model, model_opt, config):
        self.model = model
        self.model_opt = model_opt
        self.beta = config.beta
        self.dpo_epochs = config.dpo_epochs

    def train(self, inputs_ids, attention_mask, ref_logits, labels_mask):
        # 计算参考模型的logps
        ref_token_logps = self.probs_from_logits(torch.tensor(ref_logits)[:, :-1, :], inputs_ids[:, 1:])
        ref_logps = self.filter_mask(ref_token_logps, labels_mask)
        # 一次数据多训练几次，这样reference只用计算一次
        for dpo_epoch in range(self.dpo_epochs):
            # 计算策略模型的logps
            logits = self.model(inputs_ids, attention_mask)
            policy_token_logps = self.probs_from_logits(logits[:, :-1, :], inputs_ids[:, 1:])
            policy_logps = self.filter_mask(policy_token_logps, labels_mask)
            loss = self.dpo_loss(policy_logps, ref_logps)
            self.model_opt.zero_grad()
            loss.backward()
            self.model_opt.step()
            print(loss)

    def dpo_loss(self, policy_logps, ref_logps):
        """
        计算公式L_{DPO}(\pi_{\theta};\pi_{ref}) = -E[log sigmoid(\beta[log(\pi_\theta(y_w|x)/\pi_{ref}(y_w|x) -  log(\pi_\theta(y_l|x)/\pi_{ref}(y_l|x))])]
        :param policy_logps: 策略模型的logps
        :param ref_logps: 参考模型的logps
        :return: loss
        """

        def concat_probs(logps):
            """
            拆开合理与不合理数据的logps
            :param logps: 参考模型或者策略模型的logps
            :return: 合理和不合理数据的logps
            """
            len_chosen = int(len(logps) / 2)
            rejected_data = torch.cat(logps[:len_chosen])
            chosen_data = torch.cat(logps[len_chosen:])
            return rejected_data, chosen_data

        policy_rejected_logps, policy_chosen_logps = concat_probs(policy_logps)  # 计算合理数据的logps和不合理数据的logps
        ref_rejected_logps, ref_chosen_logps = concat_probs(ref_logps)
        pi_logratios = policy_chosen_logps - policy_rejected_logps
        ref_logratios = ref_chosen_logps - ref_rejected_logps
        # 这里计算策略模型合理与不合理差距与参考模型合理与不合理差距的差
        logits = pi_logratios - ref_logratios
        loss = -F.logsigmoid(self.beta * logits)
        return loss.mean()

    @staticmethod
    def probs_from_logits(logits, labels):
        log_probs = F.log_softmax(logits, dim=2)
        probs = torch.gather(log_probs, 2, labels.unsqueeze(2)).squeeze(-1)
        return probs

    @staticmethod
    def filter_mask(values, labels_masks):
        """
        :param values: 一般是prob_old、prob_ref、value(价值)的值
        :param labels_masks:label 对应的mask
        :return: 去除padding之后的数据
        """
        return [value[one_response_ids_mask[:-1] == 1].sum().unsqueeze(0) for value, one_response_ids_mask in
                zip(values, labels_masks)]
