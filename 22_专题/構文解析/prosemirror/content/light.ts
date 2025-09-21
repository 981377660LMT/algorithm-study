/**
 * @file lightweight_content_matcher.ts
 * @description A lightweight, dependency-free, programmatic content matcher.
 * It validates sequences of objects against a defined pattern without a compilation step.
 */

/**
 * Describes the objects that can be part of a sequence.
 */
export interface MatchableType {
  /** A unique name for this type. */
  readonly name: string
  /** Checks if this type belongs to a named group. */
  isInGroup(groupName: string): boolean
}

/**
 * Defines a single rule in a pattern.
 */
export interface MatchRule {
  /** The name of the type to match. */
  name?: string
  /** The name of the group to match. Any type in the group is valid. */
  group?: string
  /**
   * Defines how many times the type/group can be matched.
   * `?`: Zero or one time.
   * `*`: Zero or more times.
   * `+`: One or more times.
   * (Default): Exactly one time.
   */
  quantifier?: '?' | '*' | '+'
  /**
   * A function to create a default instance of a type.
   * Required for `fillBefore` to work on this rule.
   */
  createDefault?: () => { type: MatchableType }
}

export class LightweightContentMatcher {
  /**
   * @param pattern The array of rules defining the match pattern.
   * @param index The current position in the pattern.
   * @param counts Tracks repetitions for `*` and `+` quantifiers.
   */
  private constructor(
    readonly pattern: readonly MatchRule[],
    readonly index: number = 0,
    readonly counts: readonly number[] = []
  ) {}

  /**
   * Creates the initial matcher state from a pattern.
   * @param pattern The array of rules.
   */
  static create(pattern: readonly MatchRule[]): LightweightContentMatcher {
    return new LightweightContentMatcher(pattern, 0, new Array(pattern.length).fill(0))
  }

  /**
   * True if the current state is a valid end for the sequence.
   */
  get validEnd(): boolean {
    const rule = this.pattern[this.index]
    if (!rule) return true // End of pattern
    // Valid if the rest of the pattern is optional
    for (let i = this.index; i < this.pattern.length; i++) {
      const q = this.pattern[i].quantifier
      if (!(q === '?' || q === '*')) return false
    }
    return true
  }

  /**
   * Tries to match a type and returns the next state if successful.
   * @param type The `MatchableType` of the object to match.
   */
  matchType(type: MatchableType): LightweightContentMatcher | null {
    for (let i = this.index; i < this.pattern.length; i++) {
      const rule = this.pattern[i]
      const count = this.counts[i] || 0

      // Check if the type matches the current rule
      const matches = rule.name ? type.name === rule.name : type.isInGroup(rule.group!)

      if (matches) {
        const newCounts = [...this.counts]
        newCounts[i] = count + 1
        let nextIndex = i
        // Advance index if not a multi-match quantifier
        if (rule.quantifier !== '*' && rule.quantifier !== '+') {
          nextIndex = i + 1
        }
        return new LightweightContentMatcher(this.pattern, nextIndex, newCounts)
      }

      // If no match, check if the current rule could be skipped
      const q = rule.quantifier
      if (q === '?' || q === '*') {
        continue // Try next rule
      }
      if (q === '+' && count > 0) {
        continue // Satisfied `+`, try next rule
      }

      // Required rule was not matched
      return null
    }
    return null // Ran out of pattern to match
  }

  /**
   * Finds a sequence of default-creatable objects to insert before a given
   * `type` to make it match.
   * @param type The target type we want to be able to insert.
   */
  fillBefore(type: MatchableType): { fragment: { type: MatchableType }[] } | null {
    const fragment: { type: MatchableType }[] = []
    let current: LightweightContentMatcher = this

    for (let i = 0; i < this.pattern.length; i++) {
      // Limit search to avoid infinite loops
      if (current.matchType(type)) {
        return { fragment }
      }

      const rule = current.pattern[current.index]
      if (!rule || !rule.createDefault) return null // Cannot fill

      const filled = rule.createDefault()
      fragment.push(filled)

      const next = current.matchType(filled.type)
      if (!next) return null // Filling failed
      current = next
    }
    return null
  }
}

// --- 1. 定义你的类型 ---
const HeadingType: MatchableType = {
  name: 'heading',
  isInGroup: g => g === 'block'
}
const ParagraphType: MatchableType = {
  name: 'paragraph',
  isInGroup: g => g === 'block'
}

// --- 2. 定义你的模式 ---
const docPattern: readonly MatchRule[] = [
  {
    name: 'heading',
    // createDefault is needed for fillBefore
    createDefault: () => ({ type: HeadingType })
  },
  {
    group: 'block',
    quantifier: '*', // Zero or more blocks
    createDefault: () => ({ type: ParagraphType })
  }
]

// --- 3. 创建初始匹配器 ---
const matcher = LightweightContentMatcher.create(docPattern)

function validate(sequence: MatchableType[]): boolean {
  let current: LightweightContentMatcher | null = matcher
  for (const item of sequence) {
    if (!current) return false
    current = current.matchType(item)
  }
  return !!current && current.validEnd
}

console.log('Valid sequence:', validate([HeadingType, ParagraphType])) // true
console.log('Invalid start:', validate([ParagraphType])) // false
console.log('Valid empty blocks:', validate([HeadingType])) // true, because block* is optional
console.log('Invalid end:', validate([HeadingType, ParagraphType, HeadingType])) // false

// 在一个空文档的开头，我们想插入一个 paragraph。
// 这不合法，因为它需要一个 heading 在前面。
const initialMatcher = LightweightContentMatcher.create(docPattern)

const fill = initialMatcher.fillBefore(ParagraphType)
if (fill) {
  console.log('To insert a paragraph, first insert:')
  console.log(fill.fragment.map(f => f.type.name).join(', ')) // "heading"
} else {
  console.log('Cannot find a way to insert paragraph.')
}
