/* eslint-disable @typescript-eslint/no-empty-function */

import { DirectedGraph } from '../directed-graph/directed-graph'
import type { MultiOptionsFn, Options, OptionsFn, Runnable, Schedule, SingleOptionsFn, Tag } from './types.ts'

/**
 * Splits the input ids into tags and runnables based on their type and retrieves corresponding tags and runnables from the schedule.
 *
 * @param schedule - The schedule to retrieve tags and runnables from.
 * @param ids - The ids to split into tags and runnables.
 * @return An object containing the extracted tags and runnables.
 */
function splitTagsAndRunnables<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, ...ids: (symbol | string | Runnable<T>)[]) {
  let tags: Tag<T>[] = []
  let runnables: Runnable<T>[] = []

  for (let i = 0; i < ids.length; i++) {
    const id = ids[i]

    if (typeof id === 'symbol' || typeof id === 'string') {
      const tag = getTag(schedule, id)

      if (tag) {
        tags.push(tag)
        continue
      }

      const runnable = getRunnable(schedule, id)

      if (runnable) {
        runnables.push(runnable)
      }

      continue
    }

    runnables.push(id)
  }

  return { tags, runnables }
}

/**
 * An options function that schedules runnables before specified tags and runnables.
 *
 * @param ids - The ids to split into tags and runnables.
 * @return A function to schedule runnables before specified tags and runnables.
 */
export function before<T extends Scheduler.Context = Scheduler.Context>(...ids: (symbol | string | Runnable<T>)[]): MultiOptionsFn<T> {
  const fn: MultiOptionsFn<T> = ({ schedule, dag, runnable, tag }) => {
    const { tags, runnables } = splitTagsAndRunnables(schedule, ...ids)

    if (runnable) {
      // schedule runnable before any tags in ids
      for (const t of tags) {
        dag.addEdge(runnable, t.before)
      }

      // schedule runnable before any runnables in ids
      for (const r of runnables) {
        dag.addEdge(runnable, r)
      }
    }

    if (tag) {
      // schedule tag before any other tags in ids
      for (const t of tags) {
        dag.addEdge(tag.after, t.before)
      }

      // schedule tag before any runnables in ids
      for (const r of runnables) {
        dag.addEdge(tag.after, r)
      }
    }
  }

  fn.__type = 'multi'

  return fn
}

/**
 * An options function that schedules runnables after specified tags and runnables.
 *
 * @param ids - The ids to split into tags and runnables.
 * @return A function to schedule runnables after specified tags and runnables.
 */
export function after<T extends Scheduler.Context = Scheduler.Context>(...ids: (symbol | string | Runnable<T>)[]): MultiOptionsFn<T> {
  const fn: MultiOptionsFn<T> = ({ schedule, dag, runnable, tag }) => {
    const { tags, runnables } = splitTagsAndRunnables(schedule, ...ids)

    if (runnable) {
      // schedule runnable after any tags in ids
      for (const t of tags) {
        dag.addEdge(t.after, runnable)
      }

      // schedule runnable after any runnables in ids
      for (const r of runnables) {
        dag.addEdge(r, runnable)
      }
    }

    if (tag) {
      // schedule tag after any other tags in ids
      for (const t of tags) {
        dag.addEdge(t.after, tag.before)
      }

      // schedule tag after any runnables in ids
      for (const r of runnables) {
        dag.addEdge(r, tag.before)
      }
    }
  }

  fn.__type = 'multi'

  return fn
}

/**
 * An options function that sets the ID of a runnable in the schedule.
 *
 * @param id - The ID to be set for the runnable.
 * @return A function that sets the ID of a runnable in the schedule.
 */
export function id<T extends Scheduler.Context = Scheduler.Context>(id: symbol | string): SingleOptionsFn<T> {
  const fn: SingleOptionsFn<T> = ({ runnable, dag, schedule }) => {
    if (!runnable) {
      throw new Error('Id can only be applied to a runnable')
    }

    if (schedule.symbols.has(id)) {
      throw new Error(`Could not set id ${String(id)} because it already exists in the schedule`)
    }

    if (typeof id === 'string') {
      dag.name(runnable, id)
    }

    schedule.symbols.set(id, runnable)
  }

  fn.__type = 'single'

  return fn
}

/**
 * An options function that applies a tag to a runnable based on the given id, symbol or runnable.
 *
 * @param id - The unique identifier of the tag to apply.
 * @return A function that applies the tag to the provided runnable.
 */
export function tag<T extends Scheduler.Context = Scheduler.Context>(id: symbol | string): MultiOptionsFn<T> {
  const fn: MultiOptionsFn<T> = ({ schedule, runnable, dag }) => {
    if (!runnable) {
      throw new Error('Tag can only be applied to a runnable')
    }

    // apply the tag
    let tag = getTag(schedule, id)

    if (!tag) {
      throw new Error(`Could not find tag with id ${String(id)}`)
    }

    dag.addEdge(tag.before, runnable)
    dag.addEdge(runnable, tag.after)
  }

  fn.__type = 'multi'

  return fn
}

/**
 * Creates a new Schedule object with an empty DirectedGraph,
 * a new Map for tags, and a new Map for symbols.
 *
 * @return The newly created Schedule object.
 */
export function create<T extends Scheduler.Context = Scheduler.Context>(): Schedule<T> {
  const schedule: Schedule<T> = {
    dag: new DirectedGraph<Runnable<T>>(),
    tags: new Map(),
    symbols: new Map()
  }

  return schedule
}

/**
 * Executes all the runnables in the given schedule with the provided context.
 *
 * @param schedule - The schedule containing the runnables to execute.
 * @param context - The context to be passed to each runnable.
 */
export async function run<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, context: T) {
  for (let i = 0; i < schedule.dag.sorted.length; i++) {
    const runnable = schedule.dag.sorted[i]
    const result = runnable(context)
    if (result instanceof Promise) {
      // eslint-disable-next-line no-await-in-loop
      await result
    }
  }
}

/**
 * Removes a tag from the given schedule by its ID.
 *
 * @param schedule - The schedule from which to remove the tag.
 * @param id - The ID of the tag to remove.
 * @return This function does not return anything.
 */
export function removeTag<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, id: symbol | string) {
  const tag = schedule.tags.get(id)

  if (!tag) {
    return
  }

  schedule.dag.removeVertex(tag.before)
  schedule.dag.removeVertex(tag.after)

  schedule.tags.delete(id)
}

/**
 * Checks if a tag with the given ID exists in the schedule.
 *
 * @param schedule - The schedule to check.
 * @param id - The ID of the tag to check.
 * @return Returns true if the tag exists, false otherwise.
 */
export function hasTag<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, id: symbol | string) {
  return schedule.tags.has(id)
}

/**
 * Creates a new tag for the given schedule with the provided ID, name, and options.
 *
 * @param schedule - The schedule to create the tag for.
 * @param id - The unique identifier for the tag.
 * @param name - The name of the tag.
 * @param options - Additional options to customize the tag.
 * @return The newly created tag.
 */
export function createTag<T extends Scheduler.Context = Scheduler.Context>(
  schedule: Schedule<T>,
  id: symbol | string,
  ...options: OptionsFn<T>[]
): Tag {
  if (hasTag<T>(schedule, id)) {
    throw new Error(`Tag with id ${String(id)} already exists`)
  }

  const before = () => {}
  const after = () => {}

  const name = typeof id === 'string' ? id : String(id)

  schedule.dag.addVertex(before, {
    name: `${name}-before`,
    excludeFromSort: true
  })

  schedule.dag.addVertex(after, {
    name: `${name}-after`,
    excludeFromSort: true
  })

  schedule.dag.addEdge(before, after)

  const tag = { id, before, after }

  const optionParams: Options<T> = {
    dag: schedule.dag,
    tag,
    schedule
  }

  // apply all options: tag, before, after
  for (const option of options) {
    option(optionParams)
  }

  schedule.tags.set(id, tag)

  return tag
}

/**
 * Adds a runnable to the schedule and applies options to it.  The schedule must be built after a runnable is added.
 *
 * @param schedule - The schedule to add the runnable to.
 * @param runnable - The runnable or array of runnables to add to the schedule.
 * @param options - The options to apply to the runnable. ID can not be used if runnable is an array.
 * @throws If the runnable already exists in the schedule.
 * @return
 */
export function add<T extends Scheduler.Context = Scheduler.Context, R extends Runnable<T> | Runnable<T>[] = Runnable<T> | Runnable<T>[]>(
  schedule: Schedule<T>,
  runnable: R,
  ...options: R extends Runnable<T> ? SingleOptionsFn<T>[] : MultiOptionsFn<T>[]
) {
  let runnables: Runnable<T>[] = []

  if (Array.isArray(runnable)) {
    runnables = runnable
  } else {
    runnables = [runnable]
  }

  for (const r of runnables) {
    if (schedule.dag.exists(r)) {
      throw new Error('Runnable already exists in schedule')
    }

    // add the runnable to the graph
    schedule.dag.addVertex(r, {})

    const optionParams: Options<T> = {
      dag: schedule.dag,
      runnable: r,
      schedule
    }

    // apply all options: tag, before, after
    for (const option of options) {
      option(optionParams)
    }
  }
}

/**
 * Checks if the given runnable exists in the schedule.
 *
 * @param schedule - The schedule to check.
 * @param runnable - The runnable to check.
 * @return Returns true if the runnable exists, false otherwise.
 */
export function has<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, runnable: Runnable<T>) {
  return schedule.dag.exists(runnable)
}

/**
 * Builds the schedule by performing a topological sort on the directed graph.
 *
 * @param schedule - The schedule to be built.
 * @return This function does not return anything.
 */
export function build<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>) {
  schedule.dag.topSort()
}

/**
 * Removes a runnable from the given schedule.
 *
 * @param schedule - The schedule from which to remove the runnable.
 * @param runnable - The runnable to remove from the schedule.
 * @return This function does not return anything.
 */
export function remove<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, runnable: Runnable<T>) {
  schedule.dag.removeVertex(runnable)
}

/**
 * Retrieves a runnable from the schedule based on the given ID.
 *
 * @param schedule - The schedule to retrieve the runnable from.
 * @param id - The ID of the runnable to retrieve.
 * @return The retrieved runnable or undefined if not found.
 */
export function getRunnable<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, id: symbol | string) {
  return schedule.symbols.get(id)
}

/**
 * Retrieves a tag from the given schedule based on the provided ID.
 *
 * @param schedule - The schedule to retrieve the tag from.
 * @param id - The ID of the tag to retrieve.
 * @return The retrieved tag or undefined if not found.
 */
export function getTag<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>, id: symbol | string) {
  return schedule.tags.get(id)
}

/**
 * Prints an ASCII visualization of the directed acyclic graph (DAG) in the given schedule.
 *
 * @param schedule - The schedule containing the DAG to visualize.
 * @return This function does not return anything.
 */
export function debug<T extends Scheduler.Context = Scheduler.Context>(schedule: Schedule<T>) {
  schedule.dag.asciiVisualize()
}
