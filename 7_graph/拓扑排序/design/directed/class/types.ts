import { Runnable } from '../scheduler/types'

export type OptionsObject<T extends Scheduler.Context = Scheduler.Context> = {
  before?: symbol | string | Runnable<T> | (symbol | string | Runnable<T>)[]
  after?: symbol | string | Runnable<T> | (symbol | string | Runnable<T>)[]
  tag?: symbol | string | (symbol | string)[]
}

export type SingleOptionsObject<T extends Scheduler.Context = Scheduler.Context> = OptionsObject<T> & { id?: string | symbol }
