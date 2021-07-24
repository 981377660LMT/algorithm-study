from collections import namedtuple


class WeightedEdge(namedtuple('WeightedEdge', ['v', 'w', 'weight'])):

    def __str__(self):
        return '({}-{}: {})'.format(self.v, self.w, self.weight)

    def __repr__(self):
        return self.__str__()