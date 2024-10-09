import { type Interpreter } from '../Interpreter'

export abstract class LoxCallable {
  abstract call(interpreter: Interpreter, args: unknown[]): unknown
  abstract arity(): number
}
