interface Modify {
  index(): number
}

interface Query {
  left(): number
  right(): number
}

interface VersionQuery extends Query {
  version(): number
}

interface State<Q extends Query> {
  add(i: number): void
  remove(i: number): void
  answer(q: Q): void
}

interface AddOnlyState<Q extends Query> extends State<Q> {
  save(): void
  rollback(): void
}

interface RemoveOnlyState<Q extends Query> extends State<Q> {
  save(): void
  rollback(): void
}

interface ModifiableState<Q extends VersionQuery, M extends Modify> extends State<Q> {
  apply(m: M): void
  revoke(m: M): void
}

export {}
