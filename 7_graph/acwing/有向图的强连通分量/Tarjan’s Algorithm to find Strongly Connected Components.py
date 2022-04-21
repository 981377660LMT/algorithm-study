# https://www.geeksforgeeks.org/tarjan-algorithm-find-strongly-connected-components/?ref=lbp

# Python program to find strongly connected components in a given
# directed graph using Tarjan's algorithm (single DFS)
# Complexity : O(V+E)

from collections import defaultdict

# This class represents an directed graph
# using adjacency list representation
class Graph:
    def __init__(self, vertices):
        # No. of vertices
        self.V = vertices

        # default dictionary to store graph
        self.graph = defaultdict(list)

        self.Time = 0

    # function to add an edge to graph
    def addEdge(self, u, v):
        self.graph[u].append(v)

    '''A recursive function that find finds and prints strongly connected
    components using DFS traversal
    u --> The vertex to be visited next
    disc[] --> Stores discovery times of visited vertices
    low[] -- >> earliest visited vertex (the vertex with minimum
                discovery time) that can be reached from subtree
                rooted with current vertex
     st -- >> To store all the connected ancestors (could be part
           of SCC)
     stackMember[] --> bit/index array for faster check whether
                  a node is in stack
    '''

    def SCCUtil(self, u, low, disc, stackMember, st):

        # Initialize discovery time and low value
        disc[u] = self.Time
        low[u] = self.Time
        self.Time += 1
        stackMember[u] = True
        st.append(u)

        # Go through all vertices adjacent to this
        for v in self.graph[u]:

            # If v is not visited yet, then recur for it
            if disc[v] == -1:

                self.SCCUtil(v, low, disc, stackMember, st)

                # Check if the subtree rooted with v has a connection to
                # one of the ancestors of u
                # Case 1 (per above discussion on Disc and Low value)
                low[u] = min(low[u], low[v])

            elif stackMember[v] == True:

                '''Update low value of 'u' only if 'v' is still in stack
                (i.e. it's a back edge, not cross edge).
                Case 2 (per above discussion on Disc and Low value) '''
                low[u] = min(low[u], disc[v])

        # head node found, pop the stack and print an SCC
        w = -1  # To store stack extracted vertices
        if low[u] == disc[u]:
            while w != u:
                w = st.pop()
                print(w, end=" ")
                stackMember[w] = False

            print()

    # The function to do DFS traversal.
    # It uses recursive SCCUtil()
    def SCC(self):

        # Mark all the vertices as not visited
        # and Initialize parent and visited,
        # and ap(articulation point) arrays
        disc = [-1] * (self.V)
        low = [-1] * (self.V)
        stackMember = [False] * (self.V)
        st = []

        # Call the recursive helper function
        # to find articulation points
        # in DFS tree rooted with vertex 'i'
        for i in range(self.V):
            if disc[i] == -1:
                self.SCCUtil(i, low, disc, stackMember, st)


# Create a graph given in the above diagram
g1 = Graph(5)
g1.addEdge(1, 0)
g1.addEdge(0, 2)
g1.addEdge(2, 1)
g1.addEdge(0, 3)
g1.addEdge(3, 4)
print("SSC in first graph ")
g1.SCC()

# g2 = Graph(4)
# g2.addEdge(0, 1)
# g2.addEdge(1, 2)
# g2.addEdge(2, 3)
# print("nSSC in second graph ")
# g2.SCC()


# g3 = Graph(7)
# g3.addEdge(0, 1)
# g3.addEdge(1, 2)
# g3.addEdge(2, 0)
# g3.addEdge(1, 3)
# g3.addEdge(1, 4)
# g3.addEdge(1, 6)
# g3.addEdge(3, 5)
# g3.addEdge(4, 5)
# print("nSSC in third graph ")
# g3.SCC()

# g4 = Graph(11)
# g4.addEdge(0, 1)
# g4.addEdge(0, 3)
# g4.addEdge(1, 2)
# g4.addEdge(1, 4)
# g4.addEdge(2, 0)
# g4.addEdge(2, 6)
# g4.addEdge(3, 2)
# g4.addEdge(4, 5)
# g4.addEdge(4, 6)
# g4.addEdge(5, 6)
# g4.addEdge(5, 7)
# g4.addEdge(5, 8)
# g4.addEdge(5, 9)
# g4.addEdge(6, 4)
# g4.addEdge(7, 9)
# g4.addEdge(8, 9)
# g4.addEdge(9, 8)
# print("nSSC in fourth graph ")
# g4.SCC()


# g5 = Graph(5)
# g5.addEdge(0, 1)
# g5.addEdge(1, 2)
# g5.addEdge(2, 3)
# g5.addEdge(2, 4)
# g5.addEdge(3, 0)
# g5.addEdge(4, 2)
# print("nSSC in fifth graph ")
# g5.SCC()

# This code is contributed by Neelam Yadav

