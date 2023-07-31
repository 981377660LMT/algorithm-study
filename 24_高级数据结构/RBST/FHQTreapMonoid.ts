// !幺半群上的FHQTreap, 不支持持久化
// !因为修改会影响原树, 所以所有使用分裂合并的api都会返回新结点

interface Operation<E, Id> {}

/**
 * Proxy for a FHQTreapMonoid Node.
 */
class FHQTreapMonoid<E, Id> {}

export { FHQTreapMonoid }

if (require.main === module) {
}
