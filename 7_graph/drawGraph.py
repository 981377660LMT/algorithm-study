from typing import List, Tuple
import matplotlib.pyplot as plt
import networkx as nx


def drawGraph(edges: List[Tuple[int, int, int]], isDirected: bool) -> None:
    graph = nx.DiGraph() if isDirected else nx.Graph()
    for u, v, w in edges:
        graph.add_edge(u, v, weight=w)
    edgeLabels = {(u, v): w["weight"] for u, v, w in graph.edges(data=True)}
    pos = nx.spring_layout(graph)
    nx.draw_networkx(
        graph,
        pos,
        with_labels=True,
        node_size=1000,
        vmin=0,
        vmax=2.0,
        node_color="r",
        edge_color="r",
        font_size=20,
        font_weight="bold",
        font_family="microsoft yahei",
        alpha=0.5,
        width=2.0,
    )
    nx.draw_networkx_edge_labels(graph, pos, edge_labels=edgeLabels, font_size=15)

    plt.axis("off")
    plt.show()


if __name__ == "__main__":
    drawGraph([(1, 2, 1), (2, 3, 2), (3, 4, 3), (4, 5, 4), (1, 5, 5)], True)
    drawGraph([(1, 2, 1), (2, 3, 2), (3, 4, 3), (4, 5, 4), (1, 5, 5)], False)
