# 跳表(skiplist)

https://oi-wiki.org/ds/skiplist/

## 引言

跳表（Skip List）是一种由 William Pugh 在 1990 年提出的**概率型**数据结构，用于在有序元素集合中实现高效的查找、插入和删除操作。跳表通过在基础的有序链表上建立多级索引，优化了数据访问效率，提供了与平衡二叉查找树（如 AVL 树、红黑树）相媲美的性能，但实现更为简单。跳表因其简洁性、可扩展性以及高效的并发支持，广泛应用于数据库、内存缓存系统等领域。

本文将深入探讨跳表的结构、操作、性能分析、与其他数据结构的比较、变种与优化以及实际应用场景，并通过示例代码进一步理解其实现细节。

## 一、跳表的基本结构

### 1.1 定义

跳表是一种多层链表结构，支持高效的有序数据操作。它在基础的有序单链表上，通过构建多级索引，使得搜索速度接近二分查找，达到**O(log n)**的时间复杂度。

### 1.2 结构组成

跳表由多个层级组成，每一层都是一个**有序单链表**，且高层级的链表为低层级链表的子集。最低层（层级 1）包含所有元素，每上升一层，链表的元素数量约减少一半。具体结构如下：

- **节点（Node）**：每个节点包含一个关键字（Key）、指向下一个节点的指针数组（多个层级）。
- **层级（Level）**：跳表的层级数目，决定了索引的层数。通常基于概率决定每个节点的层级。
- **头节点（Head）**：特殊的起始节点，具有所有层级的指针，指向各层的第一个实际节点。

### 1.3 层级的确定

跳表的层级通常采用**随机化**方法决定。每个新插入的节点，通过**抛硬币**（或其他随机机制）决定其层数，概率为**1/2**。因此，每增加一层，节点数目约减半。这样的设计保证了跳表在平均情况下具有对数级别的层数。

## 二、跳表的基本操作

### 2.1 查找操作

**目标**：在跳表中查找一个特定的关键字是否存在，并返回其位置。

**步骤**：

1. 从跳表的最高层开始，设置一个指针`current`指向头节点。
2. 在当前层中，向前移动指针，直到下一个节点的关键字大于目标关键字。
3. 切换到下一层，重复步骤 2，直到最低层。
4. 在最低层进行精确查找，确定节点是否存在。

**示意图**：

```
层数: 3
索引链表:
层3: Head -> [10] -> [30] -> [50]
层2: Head -> [10] -> [30] -> [50]
层1: Head -> [10] -> [20] -> [30] -> [40] -> [50] -> [60]
```

查找`30`：

- 在层 3，从`Head`移向`10`，再移向`30`。
- 找到目标`30`，查找结束。

### 2.2 插入操作

**目标**：在跳表中插入一个新的关键字。

**步骤**：

1. **查找位置**：与查找操作类似，遍历跳表以找到插入位置的前驱节点。
2. **决定层级**：随机确定新节点的层数，通常使用概率`1/2`。
3. **插入节点**：
   - 对于新节点的每一层，将其插入到相应层的前驱节点之后。
   - 更新前驱节点的指针，使其指向新节点。
4. **更新层级**：如果新节点的层数超过当前跳表的最大层数，需增加跳表的层数，更新头节点的指针。

**代码示例（C++）：**

```cpp
#include <iostream>
#include <vector>
#include <cstdlib>
#include <ctime>
using namespace std;

struct Node {
    int key;
    vector<Node*> forward;
    Node(int k, int level) : key(k), forward(level, nullptr) {}
};

class SkipList {
private:
    int MAX_LEVEL;
    float P;
    int level;
    Node* header;

public:
    SkipList(int max_level, float p) : MAX_LEVEL(max_level), P(p), level(1) {
        header = new Node(-1, MAX_LEVEL);
    }

    ~SkipList() {
        delete header;
    }

    int randomLevel() {
        int lvl = 1;
        while (((float)rand() / RAND_MAX) < P && lvl < MAX_LEVEL)
            lvl++;
        return lvl;
    }

    void insert(int key) {
        vector<Node*> update(MAX_LEVEL, nullptr);
        Node* current = header;

        // 从高层向下查找前驱节点
        for(int i = level-1; i >=0; i--){
            while(current->forward[i] != nullptr && current->forward[i]->key < key)
                current = current->forward[i];
            update[i] = current;
        }

        current = current->forward[0];

        // 如果键不存在，插入新节点
        if(current == nullptr || current->key != key){
            int rlevel = randomLevel();

            if(rlevel > level){
                for(int i = level; i < rlevel; i++)
                    update[i] = header;
                level = rlevel;
            }

            Node* newNode = new Node(key, rlevel);
            for(int i =0; i < rlevel; i++){
                newNode->forward[i] = update[i]->forward[i];
                update[i]->forward[i] = newNode;
            }
            cout << "成功插入键 " << key << "，层数：" << rlevel << endl;
        }
    }

    // 其他操作如查找和删除...
};

int main(){
    srand((unsigned)time(0));
    SkipList list(6, 0.5);
    list.insert(3);
    list.insert(6);
    list.insert(7);
    list.insert(9);
    list.insert(12);
    list.insert(19);
    list.insert(17);
    list.insert(26);
    list.insert(21);
    list.insert(25);
    return 0;
}
```

### 2.3 删除操作

**目标**：从跳表中删除一个特定的关键字。

**步骤**：

1. **查找节点**：与查找操作类似，遍历跳表找到目标节点。
2. **更新指针**：
   - 对于每一层，更新前驱节点的指针，跳过被删除的节点。
3. **释放节点**：释放被删除节点的内存。
4. **调整层级**：如果删除后某些层级不再有节点，减少跳表的层级数。

**代码示例（C++）：**

```cpp
void deleteNode(int key){
    vector<Node*> update(MAX_LEVEL, nullptr);
    Node* current = header;

    for(int i = level-1; i >=0; i--){
        while(current->forward[i] != nullptr && current->forward[i]->key < key)
            current = current->forward[i];
        update[i] = current;
    }

    current = current->forward[0];

    if(current != nullptr && current->key == key){
        for(int i =0; i < level; i++){
            if(update[i]->forward[i] != current)
                break;
            update[i]->forward[i] = current->forward[i];
        }
        delete current;

        while(level >1 && header->forward[level-1] == nullptr)
            level--;
        cout << "成功删除键 " << key << endl;
    }
}
```

## 三、跳表的性能分析

### 3.1 时间复杂度

- **平均情况下**：
  - **查找、插入、删除**：**O(log n)**。由于跳表的层数是对数级别，通过多级索引快速定位目标节点，使得操作效率接近于平衡二叉查找树。
- **最坏情况下**：
  - **查找、插入、删除**：**O(n)**。如果所有节点都在最底层，跳表退化为单链表。但概率上，这种情况极其罕见。

### 3.2 空间复杂度

跳表需要额外的空间来存储多级指针。平均情况下，每个节点有 2 层指针，因此空间复杂度为**O(n)**。相比平衡二叉查找树，跳表在空间利用上更为高效。

### 3.3 平衡性与概率性分析

跳表通过随机化层级的方式，实现概率性的平衡。每个节点拥有多层指针的概率是指数递减的，保证了跳表整体的高度保持在**O(log n)**。这种概率性平衡的设计，使跳表在平均情况下具备优良的性能。

## 四、跳表与其他数据结构的比较

### 4.1 跳表 vs 平衡二叉查找树

| 特性           | 跳表                               | 平衡二叉查找树                           |
| -------------- | ---------------------------------- | ---------------------------------------- |
| **实现复杂度** | 相对简单，主要依赖随机化层级       | 复杂，需要维护树的平衡性（如旋转操作）   |
| **性能**       | 平均**O(log n)**，最坏**O(n)**     | 平均和最坏均为**O(log n)**               |
| **并发支持**   | 易于实现高效的并发控制（如分层锁） | 并发控制较复杂，需要更精细的锁机制       |
| **内存局部性** | 由于多级链表，内存访问相对分散     | 树结构可能导致不好的内存局部性           |
| **空间利用率** | 较高，每个节点平均有较少的额外指针 | 空间利用效率较低，每个节点需存储父子指针 |
| **维护成本**   | 插入和删除简单，依赖随机数生成层级 | 需要频繁的平衡操作，复杂度高             |

### 4.2 跳表 vs 哈希表

| 特性           | 跳表                                   | 哈希表                             |
| -------------- | -------------------------------------- | ---------------------------------- |
| **查找性能**   | 平均**O(log n)**                       | 平均**O(1)**，最坏**O(n)**         |
| **有序性**     | 支持有序遍历和范围查询                 | 不支持有序性                       |
| **内存使用**   | 较高，需存储多级指针                   | 高效，但在负载因子高时性能下降     |
| **并发支持**   | 容易实现高效的并发操作                 | 并发控制复杂，容易产生哈希冲突     |
| **实现复杂度** | 依赖随机数实现平衡                     | 需要处理哈希函数和冲突解决机制     |
| **应用场景**   | 需要有序数据操作的场景（如数据库索引） | 需要快速查找和插入的场景（如缓存） |

## 五、跳表的变种与优化

### 5.1 并发跳表

在多线程环境下，跳表可以通过适当的锁机制或无锁设计，支持高效的并发操作。常见的优化包括：

- **分段锁**：为不同的层级或节点分配独立的锁，减少锁的粒度，提升并发性能。
- **无锁跳表**：使用原子操作和比较交换（CAS）机制，实现无锁的并发控制，进一步提升性能。

**示例**：

Java 的`ConcurrentSkipListMap`是一个线程安全的跳表实现，广泛用于 Java 并发库中。

### 5.2 优化层级结构

跳表的层级结构可以根据实际数据分布动态调整，以适应不同的应用需求。例如：

- **可调整的层级概率**：改变层级生成的概率，如使用`1/4`或`1/8`，以控制层级数目。
- **多路径跳表**：在某些层级上允许多个前向指针，进一步提高查找效率。

### 5.3 压缩与合并跳表

在处理大规模数据时，可以通过压缩和合并跳表的方式，减少内存占用和提高缓存命中率。例如：

- **压缩跳表**：移除不必要的指针或节点，保留关键索引。
- **合并跳表**：将多个跳表合并为一个，减少层级数目。

## 六、跳表的应用场景

### 6.1 数据库索引

跳表作为高效的有序数据结构，广泛应用于数据库的索引系统。例如，Redis 中的`skiplist`用于实现有序集合（sorted set）的索引。

### 6.2 内存缓存系统

跳表在内存缓存系统中用于快速查找和有序数据存储。例如，LevelDB 和 RocksDB 中利用跳表实现 MemTable 的存储。

### 6.3 网络路由和负载均衡

跳表用于快速查找和管理网络路由表，支持高效的路由更新和查询。

### 6.4 实时系统

由于跳表支持高效的并发操作，适用于实时系统中的关键数据管理，如分布式协调服务（如 Consul、Etcd）中的数据存储。

## 七、跳表的实现示例

下面以 C++为例，实现一个简单的跳表，涵盖查找、插入和删除操作。

```cpp
#include <iostream>
#include <vector>
#include <cstdlib>
#include <ctime>
#include <cmath>
using namespace std;

struct Node {
    int key;
    vector<Node*> forward;

    Node(int k, int level) : key(k), forward(level, nullptr) {}
};

class SkipList {
private:
    int MAX_LEVEL;
    float P;
    int level;
    Node* header;

public:
    SkipList(int max_level, float p) : MAX_LEVEL(max_level), P(p), level(1) {
        header = new Node(-1, MAX_LEVEL); // 头节点的key为-1，后续根据需求可修改
    }

    ~SkipList() {
        delete header;
    }

    // 随机决定新节点的层数，基于概率P
    int randomLevel() {
        int lvl = 1;
        while (((float)rand() / RAND_MAX) < P && lvl < MAX_LEVEL)
            lvl++;
        return lvl;
    }

    // 插入一个新的key
    void insert(int key) {
        vector<Node*> update(MAX_LEVEL, nullptr);
        Node* current = header;

        // 查找每一层的前驱节点
        for(int i = level-1; i >=0; i--){
            while(current->forward[i] != nullptr && current->forward[i]->key < key)
                current = current->forward[i];
            update[i] = current;
        }

        current = current->forward[0];

        // 如果key不存在，则插入新节点
        if(current == nullptr || current->key != key){
            int rlevel = randomLevel();

            if(rlevel > level){
                for(int i = level; i < rlevel; i++)
                    update[i] = header;
                level = rlevel;
            }

            Node* newNode = new Node(key, rlevel);
            for(int i =0; i < rlevel; i++){
                newNode->forward[i] = update[i]->forward[i];
                update[i]->forward[i] = newNode;
            }
            cout << "Inserted key: " << key << " with level: " << rlevel << endl;
        }
    }

    // 查找一个key
    bool search(int key){
        Node* current = header;
        for(int i = level-1; i >=0; i--){
            while(current->forward[i] != nullptr && current->forward[i]->key < key)
                current = current->forward[i];
        }
        current = current->forward[0];
        if(current != nullptr && current->key == key){
            return true;
        }
        return false;
    }

    // 删除一个key
    void deleteNode(int key){
        vector<Node*> update(MAX_LEVEL, nullptr);
        Node* current = header;

        for(int i = level-1; i >=0; i--){
            while(current->forward[i] != nullptr && current->forward[i]->key < key)
                current = current->forward[i];
            update[i] = current;
        }

        current = current->forward[0];

        if(current != nullptr && current->key == key){
            for(int i =0; i < level; i++){
                if(update[i]->forward[i] != current)
                    break;
                update[i]->forward[i] = current->forward[i];
            }
            delete current;

            while(level >1 && header->forward[level-1] == nullptr)
                level--;
            cout << "Deleted key: " << key << endl;
        }
    }

    // 打印跳表
    void display(){
        cout << "\n***** Skip List *****\n";
        for(int i=0; i < level; i++){
            Node* node = header->forward[i];
            cout << "Level " << i+1 << ": ";
            while(node != nullptr){
                cout << node->key << " ";
                node = node->forward[i];
            }
            cout << endl;
        }
    }
};

int main(){
    srand((unsigned)time(0));
    SkipList list(4, 0.5);
    list.insert(3);
    list.insert(6);
    list.insert(7);
    list.insert(9);
    list.insert(12);
    list.insert(19);
    list.insert(17);
    list.insert(26);
    list.insert(21);
    list.insert(25);
    list.display();

    // 查找操作
    int searchKey = 19;
    cout << "Searching for key " << searchKey << ": " << (list.search(searchKey) ? "Found" : "Not Found") << endl;

    // 删除操作
    list.deleteNode(19);
    list.display();
    cout << "Searching for key " << searchKey << ": " << (list.search(searchKey) ? "Found" : "Not Found") << endl;

    return 0;
}
```

**输出示例**：

```
Inserted key: 3 with level: 2
Inserted key: 6 with level: 1
Inserted key: 7 with level: 1
Inserted key: 9 with level: 3
Inserted key: 12 with level: 1
Inserted key: 19 with level: 1
Inserted key: 17 with level: 1
Inserted key: 26 with level: 2
Inserted key: 21 with level: 1
Inserted key: 25 with level: 1

***** Skip List *****
Level 1: 3 6 7 9 12 17 19 21 25 26
Level 2: 3 26
Level 3: 9

Searching for key 19: Found
Deleted key: 19

***** Skip List *****
Level 1: 3 6 7 9 12 17 21 25 26
Level 2: 3 26
Level 3: 9

Searching for key 19: Not Found
```

### 7.1 代码说明

- **Node 结构**：每个节点包含一个关键字和一个指向各层前向节点的指针数组。
- **SkipList 类**：
  - `MAX_LEVEL`：跳表的最大层级。
  - `P`：概率因子，用于决定新节点的层级（通常为 0.5）。
  - `level`：当前跳表的层级数。
  - `header`：头节点，包含所有层级的指针。
  - `randomLevel()`：根据概率 P 随机生成节点的层级。
  - `insert(int key)`：插入新节点。
  - `search(int key)`：查找是否存在某个关键字。
  - `deleteNode(int key)`：删除指定关键字的节点。
  - `display()`：打印跳表的各层结构。
- **主函数**：
  - 初始化跳表，插入若干节点，展示跳表结构。
  - 执行查找和删除操作，验证跳表功能。

## 八、跳表的优缺点

### 8.1 优点

1. **简单易实现**：相比于平衡二叉查找树，跳表的实现更为直观，依赖较少的旋转操作。
2. **高效的查找性能**：平均情况下，跳表的查找、插入和删除操作均为**O(log n)**，接近平衡树的性能。
3. **灵活的层级结构**：通过概率方式动态调整层级，适应不同的数据分布。
4. **良好的并发支持**：跳表天然适合实现并发操作，特别是通过分层锁或无锁设计。
5. **空间利用率高**：相比于红黑树等平衡树，跳表需要更少的额外空间储存指针。

### 8.2 缺点

1. **较高的空间开销**：每个节点需要存储多层指针，增加了内存使用。
2. **最坏情况退化**：在极端情况下，如所有节点都在最底层，跳表性能退化为**O(n)**。
3. **依赖随机性**：跳表的性能基于随机化层级，可能导致不均衡的层级分布，虽然概率上不常见。
4. **缓存局部性较差**：跳表的多级链表可能导致频繁的内存跳转，影响 CPU 缓存性能。

## 九、跳表的实际应用

### 9.1 Redis 的有序集合（Sorted Set）

Redis 中`sorted set`（有序集合）采用跳表作为底层数据结构，通过跳表和哈希表的组合，实现了高效的有序数据存储与快速查找。

### 9.2 LevelDB 与 RocksDB 的 MemTable

这两款高性能键值存储系统在 MemTable 的实现中，采用跳表作为内存中的有序数据结构，支持快速的插入和范围查询。

### 9.3 Java 的 ConcurrentSkipListMap

Java 的`ConcurrentSkipListMap`是一个线程安全的跳表实现，提供了高效的并发数据访问，广泛应用于 Java 并发库中。

### 9.4 其他应用

跳表在网络路由表、实时系统的数据管理、分布式锁服务（如 Consul、Etcd）等多种场景中得到应用。

## 十、结论

跳表作为一种概率型有序数据结构，通过多级链表的设计，实现了高效的查找、插入和删除操作。其简单的实现方式、良好的并发支持和接近平衡树的性能，使其在数据库索引、内存缓存系统等领域得到了广泛应用。尽管跳表存在空间开销较高和最坏情况下性能退化等缺点，但通过优化层级结构、并发设计等方法，能够有效提升其在实际应用中的表现。

通过本文的深入分析，希望读者能够全面理解跳表的原理、实现与应用，为在实际项目中选择和实现高效的数据结构提供参考。
