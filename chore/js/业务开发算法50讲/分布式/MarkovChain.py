import numpy as np


class MarkovChain:
    """
    马尔可夫链模型。
    """

    def __init__(self, transition_matrix, states):
        """
        初始化马尔可夫链。
        :param transition_matrix: 转移概率矩阵
        :param states: 状态列表
        """
        self.transition_matrix = np.array(transition_matrix)
        self.transition_matrix = self.transition_matrix / self.transition_matrix.sum(
            axis=1, keepdims=True
        )
        self.states = states
        self.state_dict = {state: i for i, state in enumerate(states)}

    def next_state(self, current_state):
        """
        根据当前状态，返回下一个状态。
        """
        current_index = self.state_dict[current_state]
        probabilities = self.transition_matrix[current_index]
        return np.random.choice(self.states, p=probabilities)

    def generate_states(self, start_state, n):
        """
        生成长度为 n 的状态序列。
        """
        current_state = start_state
        state_sequence = [current_state]
        for _ in range(n - 1):
            current_state = self.next_state(current_state)
            state_sequence.append(current_state)
        return state_sequence


if __name__ == "__main__":
    states = ["A", "B", "C"]
    transition_matrix = [[0.1, 0.6, 0.3], [0.4, 0.5, 0.1], [0.3, 0.3, 0.4]]
    mc = MarkovChain(transition_matrix, states)
    print(mc.generate_states("A", 10))
