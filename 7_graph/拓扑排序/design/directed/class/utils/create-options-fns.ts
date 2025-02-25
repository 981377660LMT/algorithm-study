import { after as afterFn, before as beforeFn, id as idFn, tag as tagFn } from '../../scheduler/scheduler'
import type { MultiOptionsFn, OptionsFn, SingleOptionsFn } from '../../scheduler/types'
import type { OptionsObject, SingleOptionsObject } from '../types'

export function createOptionsFns<T extends Scheduler.Context = Scheduler.Context>(options?: SingleOptionsObject<T>): SingleOptionsFn<T>[]
export function createOptionsFns<T extends Scheduler.Context = Scheduler.Context>(options?: OptionsObject<T>): MultiOptionsFn<T>[]
export function createOptionsFns<T extends Scheduler.Context = Scheduler.Context>(
  options?: OptionsObject<T> | SingleOptionsObject<T>
): SingleOptionsFn<T>[] | MultiOptionsFn<T>[] {
  if (!options) return []

  const optionsFns: OptionsFn<T>[] = []

  if ('id' in options && options.id) {
    optionsFns.push(idFn(options.id))
  }

  if (options.before) {
    if (Array.isArray(options.before)) {
      optionsFns.push(...options.before.map(before => beforeFn(before)))
    } else {
      optionsFns.push(beforeFn(options.before))
    }
  }

  if (options.after) {
    if (Array.isArray(options.after)) {
      optionsFns.push(...options.after.map(after => afterFn(after)))
    } else {
      optionsFns.push(afterFn(options.after))
    }
  }

  if (options.tag) {
    if (Array.isArray(options.tag)) {
      optionsFns.push(...options.tag.map(tag => tagFn<T>(tag)))
    } else {
      optionsFns.push(tagFn(options.tag))
    }
  }

  return optionsFns
}
