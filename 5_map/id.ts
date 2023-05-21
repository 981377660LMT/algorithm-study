const _pool = new Map<unknown, number>()
function id(o: unknown): number {
  if (!_pool.has(o)) {
    _pool.set(o, _pool.size)
  }
  return _pool.get(o)!
}

export {}
