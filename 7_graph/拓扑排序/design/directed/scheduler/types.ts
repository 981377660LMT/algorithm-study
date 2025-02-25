import { DirectedGraph } from '../directed-graph/types'

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Scheduler {
    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    export interface Context {}
  }
}

export type Runnable<T extends Scheduler.Context = Scheduler.Context> = (context: T) => void | Promise<void>

export interface Schedule<T extends Scheduler.Context = Scheduler.Context> {
  dag: DirectedGraph<Runnable<T>>
  tags: Map<symbol | string, Tag<T>>
  symbols: Map<symbol | string, Runnable<T>>
}

export type Options<T extends Scheduler.Context = Scheduler.Context> = {
  schedule: Schedule<T>
  dag: DirectedGraph<Runnable<T>>
  runnable?: Runnable<T>
  tag?: Tag<T>
}

export type OptionsFn<T extends Scheduler.Context = Scheduler.Context> = SingleOptionsFn<T> | MultiOptionsFn<T>

export type SingleOptionsFn<T extends Scheduler.Context = Scheduler.Context> = ((options: Options<T>) => void) & { __type: 'single' | 'multi' }

export type MultiOptionsFn<T extends Scheduler.Context = Scheduler.Context> = ((options: Options<T>) => void) & { __type: 'multi' }

export type Tag<T extends Scheduler.Context = Scheduler.Context> = {
  id: symbol | string
  before: Runnable<T>
  after: Runnable<T>
}
