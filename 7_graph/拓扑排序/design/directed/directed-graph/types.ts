export interface Vertex<T> {
  name: string
  value: T
  inEdges: Set<Vertex<T>>
  outEdges: Set<Vertex<T>>
  excludeFromSort: boolean
  get outEdgeCount(): number
  get inEdgeCount(): number
}

export interface DirectedGraph<T> {
  addVertex(value: T, options: { name?: string; excludeFromSort?: boolean }): Vertex<T>
  topSort(): T[]
  exists(value: T): boolean
  hasEdge(from: T, to: T): boolean
  addEdge(from: T, to: T): void
  excludeFromSort(value: T, exclude: boolean): void
  name(value: T, name: string): void
  removeVertex(value: T): void
  removeEdge(from: T, to: T): void
  getVertex(value: T): Vertex<T> | undefined
  addVertexToEndOfGraph(value: T, options: { name?: string; excludeFromSort?: boolean }): void
  isPathWithoutDirectEdge(from: Vertex<T>, to: Vertex<T>): boolean
  transitiveReduction(): void
  asciiVisualize(): void

  sorted: T[]
}
