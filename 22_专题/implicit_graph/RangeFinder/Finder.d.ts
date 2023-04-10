interface Finder {
  prev(x: number): number | null
  next(x: number): number | null
  erase(x: number): void
}
