// Don't do like this.
// Keep it simple and stupid.
//
// !Better: function validate(v): E[]
// Example:
// type Rule = (schema: Schema) => IErrorInfo[];
// const rules: Rule[] = [validateDuplicateComponentNames, validateDuplicateCodeNames];
// export function validate(schema: Schema): IErrorInfo[] {
//   return rules.flatMap((rule) => rule(schema));
// }

type ValidationRule<V = unknown, E = unknown> = (value: V) => E[]

class Validator<V = unknown, E = unknown> {
  private readonly _rules: ValidationRule<V, E>[] = []

  addRule(rule: ValidationRule<V, E>): void {
    this._rules.push(rule)
  }

  validate(value: V): E[] {
    const res: E[] = []
    for (const rule of this._rules) {
      const errors = rule(value)
      res.push(...errors)
    }
    return res
  }
}

export function createValidator<V = unknown, E = unknown>(
  rules: ValidationRule<V, E>[]
): Validator<V, E> {
  const res = new Validator<V, E>()
  for (const rule of rules) {
    res.addRule(rule)
  }
  return res
}

if (require.main === module) {
  const validator = createValidator([
    (value: number) => {
      if (value < 0) {
        return ['不能小于0']
      }
      return []
    },
    (value: number) => {
      if (value > 10) {
        return ['不能大于10']
      }
      return []
    },
    (value: number) => {
      if (value % 2 !== 0) {
        return ['不能是奇数']
      }
      return []
    }
  ])

  console.log(validator.validate(5))
  console.log(validator.validate(11))
  console.log(validator.validate(-1))
  console.log(validator.validate(8))
}
