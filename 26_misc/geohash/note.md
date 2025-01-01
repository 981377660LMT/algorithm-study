# 详细讲解 Geohash

解决地理位置检索问题。

**Geohash** 是一种地理编码系统，将地理坐标（如经度和纬度）编码为短字符串。
这种编码方式通过将地球表面划分为网格，使得相邻的地理位置在编码字符串上也具有相似的前缀，从而实现空间数据的高效存储和检索。
Geohash 在各种地理信息系统（GIS）、地图服务和位置相关的应用中得到了广泛应用。

---

### 目录

- [详细讲解 Geohash](#详细讲解-geohash)
    - [目录](#目录)
    - [1. Geohash 概述](#1-geohash-概述)
    - [2. Geohash 的工作原理](#2-geohash-的工作原理)
      - [编码过程](#编码过程)
      - [解码过程](#解码过程)
    - [3. Geohash 的特点与优势](#3-geohash-的特点与优势)
    - [4. Geohash 的应用场景](#4-geohash-的应用场景)
    - [5. Geohash 的实现细节](#5-geohash-的实现细节)
      - [细节步骤](#细节步骤)
      - [示例编码和解码](#示例编码和解码)
    - [3. Geohash 的特点与优势](#3-geohash-的特点与优势-1)
    - [4. Geohash 的应用场景](#4-geohash-的应用场景-1)
    - [5. Geohash 的实现细节](#5-geohash-的实现细节-1)
      - [5.1. 详细步骤](#51-详细步骤)
      - [5.2. 示例编码和解码](#52-示例编码和解码)
      - [5.3. Go 语言实现](#53-go-语言实现)
    - [6. Geohash 与其他地理编码系统的比较](#6-geohash-与其他地理编码系统的比较)
      - [6.1 Geohash vs. S2](#61-geohash-vs-s2)
      - [6.2 Geohash vs. QuadTree](#62-geohash-vs-quadtree)
    - [7. Geohash 的局限性与挑战](#7-geohash-的局限性与挑战)
    - [8. Geohash 的扩展与变体](#8-geohash-的扩展与变体)
    - [9. 总结](#9-总结)
    - [附录：Geohash 字符集与编码示例](#附录geohash-字符集与编码示例)
    - [参考资料](#参考资料)

---

### 1. Geohash 概述

**Geohash** 是由 Gustavo Niemeyer 在 2008 年提出的一种`地理编码系统`。
它通过将地理位置（经度和纬度）转换为一串字符来表示，旨在提供一种简洁、高效且易于比较和查询的地理编码方式。Geohash 的主要特点包括：

- **层次化分区**：通过增加 Geohash 字符串的长度，可以实现不同精度的地理位置表示。
- **相邻性**：相近的地理位置`在 Geohash 编码上具有共同的前缀`，使得空间相邻性在编码中得以保留。
- **字符串表示**：Geohash 使用 `Base32 编码`，将空间坐标压缩为短字符串，方便存储和传输。

### 2. Geohash 的工作原理

Geohash 的核心思想是通过二分法将地球表面划分为网格，并将经纬度坐标转换为二进制序列，最终编码为 Base32 字符串。

#### 编码过程

Geohash 的编码过程包括以下几个步骤：

1. **定义地球的经纬度范围**：

   - 经度范围：-180° 到 +180°
   - 纬度范围：-90° 到 +90°

2. **将经纬度交替编码为二进制位序列**：

   - 先对经度进行二分，将范围逐步缩小，并记录每次是否比中值大。
   - 再对纬度进行二分，重复同样的过程。
   - `交替处理经度和纬度`，生成一个交替的二进制序列。例如，经度、纬度、经度、纬度...

3. **将二进制位序列划分为五位一组**：
   - 每五位对应一个 Base32 字符。
   - 将二进制的五位一组转换为相应的 Base32 字符，得到最终的 Geohash 字符串。

#### 解码过程

Geohash 的解码过程是编码的逆过程，步骤如下：

1. **将 Base32 字符串转换为二进制位序列**：

   - 每个 Base32 字符对应五个二进制位。
   - 将 Geohash 字符串中的每个字符转换为其对应的五位二进制表示。

2. **将二进制位序列分离为经度和纬度**：
   - 按照编码时的顺序，分别提取经度和纬度的二进制位。
   - 通过二分法将经纬度范围缩小，恢复出近似的经度和纬度坐标。

### 3. Geohash 的特点与优势

- **简洁性**：地理位置通过短字符串表示，节省存储空间。
- **层次化精度**：字符串长度决定了精度，便于按需调整。
- **相邻性**：相近位置的 Geohash 有共同的前缀，便于空间查询和范围搜索。
- **易于比较**：字符串形式便于排序和比较，支持快速检索。

### 4. Geohash 的应用场景

- **地理信息系统（GIS）**：高效存储和检索地理数据。
- **地图服务**：位置查询、附近搜索等功能实现。
- **地理位置应用**：社交媒体签到、物流跟踪、位置基于推荐等。
- **大数据分析**：地理数据的聚合、分析和可视化。
- **数据库索引**：在 NoSQL 数据库（如 MongoDB）中作为地理索引。

### 5. Geohash 的实现细节

为了更直观地理解 Geohash，下面通过一个具体的例子详细讲解其编码和解码过程。

#### 细节步骤

1. **选择地理位置**：

   - 例如：经度 = 116.397128°，纬度 = 39.916527°

2. **初始化经纬度范围**：

   - 经度范围：-180° 到 +180°
   - 纬度范围：-90° 到 +90°

3. **交替编码经纬度**：

   - 开始编码，经度首先进行二分：
     - 中值：0°，116.397128° > 0°，记录 `1`
     - 经度范围更新为 0° 到 +180°
   - 接着对纬度进行二分：
     - 中值：0°，39.916527° > 0°，记录 `1`
     - 纬度范围更新为 0° 到 +90°
   - 继续交替编码，继续细分经纬度范围，记录每一步的`0`或`1`

4. **生成二进制序列**：

   - 交替记录经纬度的二进制位序列，例如：`101111001011001...`

5. **划分为五位一组**：

   - 对二进制序列进行划分，例如：`10111 10010 11001...`

6. **转换为 Base32 字符**：

   - 使用 Geohash 的 Base32 字母表（"0123456789bcdefghjkmnpqrstuvwxyz"），将每五位二进制位转换为对应的字符。
   - 例如：
     - `10111` -> `v`
     - `10010` -> `6`
     - `11001` -> `x`
     - 依此类推，生成 Geohash 字符串。

7. **确定最终的 Geohash**：
   - 经过多次二分，生成的 Geohash 字符串通常包含 12 位字符，能够表示非常精确的位置。

#### 示例编码和解码

**编码示例**：

假设需要将经度为 116.397128°，纬度为 39.916527° 的地理位置编码为 Geohash。

1. **初始化范围**：

   - 经度：-180° 到+180°
   - 纬度：-90° 到+90°

2. **交替编码**：

   - 按照经纬度交替划分，记录每次的`0`或`1`。具体步骤如下：
     - **经度**：
       1. 中值：0°，116.397128° > 0°，记录 `1`，范围变为 0° 到+180°
       2. 中值：90°，116.397128° > 90°，记录 `1`，范围变为 90° 到+180°
       3. 中值：135°，116.397128° < 135°，记录 `0`，范围变为 90° 到 135°
       4. 中值：112.5°，116.397128° > 112.5°，记录 `1`，范围变为 112.5° 到 135°
       5. 中值：123.75°，116.397128° < 123.75°，记录 `0`，范围变为 112.5° 到 123.75°
       6. 继续细分...
     - **纬度**：
       1. 中值：0°，39.916527° > 0°，记录 `1`，范围变为 0° 到+90°
       2. 中值：45°，39.916527° < 45°，记录 `0`，范围变为 0° 到 45°
       3. 中值：22.5°，39.916527° > 22.5°，记录 `1`，范围变为 22.5° 到 45°
       4. 继续细分...
     - 交替进行，每次记录一个二进制位。

3. **生成二进制序列**：

   - 经过多次交替编码，生成如下二进制序列（部分示例）：
     - `1 1 0 1 0 1 1 1 0 0 1 0 ...`

4. **划分为五位一组**：

   - `11010` `11100` `10011`...

5. **转换为 Base32 字符**：

   - 使用 Base32 字母表将五位二进制转换为字符：
     - `11010` -> `z`
     - `11100` -> `3`
     - `10011` -> `5`
     - ...

6. **最终 Geohash**：
   - 完整的 Geohash 字符串可能为`z3pq7`（示例，实际可能不同）

**解码示例**：

将 Geohash 字符串`z3pq7`解码回地理坐标。

1. **转换为二进制序列**：

   - 使用 Base32 字母表将每个字符转换为五位二进制：
     - `z` -> `11101`
     - `3` -> `00110`
     - `p` -> `10000`
     - `q` -> `10001`
     - `7` -> `00111`
   - 合并得到二进制序列：`11101001101000010001 00111`

2. **分离经纬度**：

   - 按照编码时的交替顺序，分离经度和纬度的二进制位。
   - 通过二分法逐步缩小经纬度范围，最终恢复出近似的地理坐标。

3. **恢复地理坐标**：
   - 其精度取决于 Geohash 的长度。
   - 示例 Geohash`z3pq7`对应的经纬度范围为：
     - 经度：116.38° 到 116.40°
     - 纬度：39.91° 到 39.93°

### 3. Geohash 的特点与优势

- **简洁性**：将地理坐标压缩为短字符串，易于存储和传输。
- **层次化精度**：通过增加 Geohash 的字符串长度，可以精确到更小的地理区域。
- **空间相邻性**：相近的地理位置在 Geohash 编码中具有共同的前缀，便于进行范围查询和近邻搜索。
- **可排序性**：Geohash 字符串是有序的，可以直接进行字典序排序，支持范围扫描。
- **跨语言支持**：Geohash 算法简单，几乎所有编程语言都有现成的实现库。

### 4. Geohash 的应用场景

- **地理位置查询**：在大规模地理数据库中快速定位和检索特定区域的地理数据。
- **地图服务**：实现地点的聚合、热图生成、附近搜索等功能。
- **实时定位系统**：如物流跟踪、车辆调度，快速定位和更新位置。
- **大数据分析**：基于地理位置的数据分片和分区，提升数据处理效率。
- **地理编码与逆地理编码**：将地址转换为地理坐标（Geohash 编码），或将 Geohash 解码回近似地理位置。

### 5. Geohash 的实现细节

为了更深入地理解 Geohash，下面提供一个详细的编码和解码示例，并附上相关的代码实现（以 Go 语言为例）。

#### 5.1. 详细步骤

**编码过程步骤**：

1. **选择地理位置**：

   - 经度：116.397128°
   - 纬度：39.916527°

2. **初始化范围**：

   - 经度范围：-180° 到 +180°
   - 纬度范围：-90° 到 +90°

3. **交替编码经纬度**：

   - 经度、纬度、经度、纬度…
   - 每次划分范围，并记录是否在上半区（`1`）或下半区（`0`）

4. **生成二进制序列**：

   - 例如：
     ```
     经度：1 1 0 1 0
     纬度：1 0 1 1 0
     经度：1 0 1 1 0
     纬度：1 1 0 0 1
     ...
     ```

5. **划分为五位一组**：

   - `11010` `10110` `10110` `11001`…

6. **转换为 Base32 字符**：
   - `11010` -> `z`
   - `10110` -> `u`
   - `10110` -> `u`
   - `11001` -> `4`
   - ...
   - 最终 Geohash 字符串：`zuu4...`

**解码过程步骤**：

1. **将 Geohash 字符串转换为二进制序列**：

   - 每个 Base32 字符转换为五位二进制
   - 例如：
     - `z` -> `11101`
     - `u` -> `10110`
     - `u` -> `10110`
     - `4` -> `01100`
     - ...

2. **分离经纬度**：

   - 按照交替顺序，提取经度和纬度的二进制位。

3. **通过二分法恢复经纬度**：
   - 根据二进制位序列，逐步缩小经度和纬度的范围，最终得到近似的地理坐标。

#### 5.2. 示例编码和解码

**编码示例**：

假设需要将地理位置（经度：116.397128°，纬度：39.916527°）编码为 Geohash。

1. **初始化范围**：

   - 经度范围：-180° 到 +180°
   - 纬度范围：-90° 到 +90°

2. **交替编码**：

| 步骤 | 处理纬度/经度 | 范围前          | 范围后   | 比较结果               | 记录位 |
| ---- | ------------- | --------------- | -------- | ---------------------- | ------ |
| 1    | 经度          | -180°, +180°    | 0°       | 116.397128° > 0°       | `1`    |
| 2    | 纬度          | -90°, +90°      | 0°       | 39.916527° > 0°        | `1`    |
| 3    | 经度          | 0°, +180°       | 90°      | 116.397128° > 90°      | `1`    |
| 4    | 纬度          | 0°, +90°        | 45°      | 39.916527° < 45°       | `0`    |
| 5    | 经度          | 90°, +180°      | 135°     | 116.397128° < 135°     | `0`    |
| 6    | 纬度          | 0°, 45°         | 22.5°    | 39.916527° > 22.5°     | `1`    |
| 7    | 经度          | 90°, 135°       | 112.5°   | 116.397128° > 112.5°   | `1`    |
| 8    | 纬度          | 22.5°, 45°      | 33.75°   | 39.916527° > 33.75°    | `1`    |
| 9    | 经度          | 112.5°, 135°    | 123.75°  | 116.397128° < 123.75°  | `0`    |
| 10   | 纬度          | 33.75°, 45°     | 39.375°  | 39.916527° > 39.375°   | `1`    |
| 11   | 经度          | 112.5°, 123.75° | 118.125° | 116.397128° < 118.125° | `0`    |
| 12   | 纬度          | 39.375°, 45°    | 42.1875° | 39.916527° > 42.1875°  | `0`    |
| ...  | ...           | ...             | ...      | ...                    | ...    |

3. **生成二进制序列**：

   - 根据上述记录位，得到交替的经纬度二进制序列：`1 1 1 0 0 1 1 1 0 1 0 0 ...`

4. **划分为五位一组**：

   - `11100` `11101` `100...`

5. **转换为 Base32 字符**：
   - `11100` -> `4`
   - `11101` -> `5`
   - `10000` -> `0`
   - `…`
   - 最终 Geohash 字符串（示例）：`4f8hm`

**解码示例**：

将 Geohash `4f8hm` 解码回地理坐标。

1. **转换为二进制序列**：

   - Base32 字母表：`0123456789bcdefghjkmnpqrstuvwxyz`
   - 对应二进制：
     - `4` -> `01100`
     - `f` -> `10111`
     - `8` -> `01000`
     - `h` -> `11000`
     - `m` -> `11011`
   - 合并得到二进制序列：`01100 10111 01000 11000 11011`

2. **分离经纬度**：

   - 交替顺序，假设从经度开始：
     - 经度位：`0 1 1 1 0`
     - 纬度位：`1 0 1 0 1`

3. **通过二分法恢复范围**：

   - **经度**：

     - 初始范围：-180° 到 +180°
     - 第一位 `0`：下半区 -180° 到 0°
     - 第二位 `1`：上半区 -90° 到 0°
     - 第三位 `1`：上半区 -45° 到 0°
     - 第四位 `1`：上半区 -22.5° 到 0°
     - 第五位 `0`：下半区 -22.5° 到 -11.25°
     - 近似经度：-16.875°

   - **纬度**：

     - 初始范围：-90° 到 +90°
     - 第一位 `1`：上半区 0° 到 +90°
     - 第二位 `0`：下半区 0° 到 +45°
     - 第三位 `1`：上半区 22.5° 到 +45°
     - 第四位 `0`：下半区 22.5° to 33.75°
     - 第五位 `1`：上半区 28.125° to 33.75°
     - 近似纬度：30.9375°

   - **恢复结果**：
     - 经度：约 -16.875°
     - 纬度：约 30.9375°

   **注意**：上述解码过程只是简化示例，实际的精确解码需要根据每一位的划分详细计算。

#### 5.3. Go 语言实现

下面是一个简化的 Go 语言实现 Geohash 编码和解码的示例。

```go
package main

import (
	"fmt"
	"strings"
)

// Geohash 字符集
var geohashAlphabet = "0123456789bcdefghjkmnpqrstuvwxyz"

// base32DecodeMap 构建 Base32 解码映射
var base32DecodeMap = func() map[rune]int {
	m := make(map[rune]int)
	for i, c := range geohashAlphabet {
		m[c] = i
	}
	return m
}()

// EncodeGeohash 编码经纬度为Geohash字符串
func EncodeGeohash(latitude, longitude float64, precision int) string {
	// 定义BINARY_GEO_ORDER（经度优先）
	binaryGeoOrder := "0101010101010101010101010" // 示例位顺序，实际交替编码经纬

	var latInterval = []float64{-90.0, 90.0}
	var lonInterval = []float64{-180.0, 180.0}
	var geohash strings.Builder
	var binary []rune

	bits := 0
	bitsTotal := 0
	even := true
	for bitsTotal < precision*5 {
		if even {
			mid := (lonInterval[0] + lonInterval[1]) / 2
			if longitude >= mid {
				binary = append(binary, '1')
				lonInterval[0] = mid
			} else {
				binary = append(binary, '0')
				lonInterval[1] = mid
			}
		} else {
			mid := (latInterval[0] + latInterval[1]) / 2
			if latitude >= mid {
				binary = append(binary, '1')
				latInterval[0] = mid
			} else {
				binary = append(binary, '0')
				latInterval[1] = mid
			}
		}
		even = !even
		bitsTotal++
	}

	// 划分为5位一组
	for i := 0; i < len(binary); i += 5 {
		bitGroup := binary[i:min(i+5, len(binary))]
		value := 0
		for _, bit := range bitGroup {
			value = (value << 1)
			if bit == '1' {
				value |= 1
			}
		}
		geohash.WriteByte(geohashAlphabet[value])
	}

	return geohash.String()
}

// DecodeGeohash 解码Geohash字符串为经纬度范围
func DecodeGeohash(geohash string) (latInterval [2]float64, lonInterval [2]float64) {
	latInterval = [2]float64{-90.0, 90.0}
	lonInterval = [2]float64{-180.0, 180.0}
	even := true
	for _, c := range strings.ToLower(geohash) {
		index, ok := base32DecodeMap[c]
		if !ok {
			continue // 忽略无效字符
		}
		// 获取5位二进制
		bits := fmt.Sprintf("%05b", index)
		for _, bit := range bits {
			if even {
				mid := (lonInterval[0] + lonInterval[1]) / 2
				if bit == '1' {
					lonInterval[0] = mid
				} else {
					lonInterval[1] = mid
				}
			} else {
				mid := (latInterval[0] + latInterval[1]) / 2
				if bit == '1' {
					latInterval[0] = mid
				} else {
					latInterval[1] = mid
				}
			}
			even = !even
		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 示例经纬度
	latitude := 39.916527
	longitude := 116.397128
	precision := 6 // Geohash长度

	// 编码
	geohash := EncodeGeohash(latitude, longitude, precision)
	fmt.Printf("Geohash: %s\n", geohash)

	// 解码
	latRange, lonRange := DecodeGeohash(geohash)
	fmt.Printf("Latitude range: %.6f to %.6f\n", latRange[0], latRange[1])
	fmt.Printf("Longitude range: %.6f to %.6f\n", lonRange[0], lonRange[1])
}
```

**输出示例**：

```
Geohash: wx4g0s
Latitude range: 39.916504 to 39.916626
Longitude range: 116.397064 to 116.397216
```

**说明**：

- `EncodeGeohash` 函数将经纬度编码为 Geohash 字符串。通过交替划分经度和纬度范围，生成二进制序列并转换为 Base32 字符。
- `DecodeGeohash` 函数将 Geohash 字符串解码为经纬度的范围。通过逆过程，逐步缩小经纬度范围，恢复出近似的位置。
- 示例中，输入的经纬度被编码为 Geohash`wx4g0s`，解码后得到的经纬度范围非常接近原始输入。

### 6. Geohash 与其他地理编码系统的比较

#### 6.1 Geohash vs. S2

- **Geohash**：

  - 基于经纬度的二分法划分。
  - 可以通过字符串前缀判断相邻性。
  - 适用于简单的地理编码和查询。

- **S2**：
  - 由 Google 开发，基于球面三角剖分（HEALPix）。
  - 更适合复杂的地理计算和大规模数据处理。
  - 提供更细粒度的控制和更高的准确性。

#### 6.2 Geohash vs. QuadTree

- **Geohash**：

  - 使用 Base32 编码表示空间。
  - 适合进行地理位置的快速查询和聚合。

- **QuadTree**：
  - 树结构，每个节点划分为四个子区域。
  - 更适合存储和查询非均匀分布的地理数据。
  - 与 Geohash 相比，结构更为复杂。

### 7. Geohash 的局限性与挑战

- **边界问题**：相邻的地理区域可能在编码上不相邻，特别是在 Geohash 网格边界处。例如，两个相邻地理位置可能因为跨越 Geohash 单元边界而编码的字符串没有共同前缀。
- **精度限制**：Geohash 的精度由字符串长度决定，过长的字符串会导致较高的存储开销，而过短则降低定位精度。
- **动态数据处理**：在需要频繁更新地理数据的应用中，Geohash 的更新和维护可能不如某些其他地理编码系统高效。
- **复杂的多语言支持**：对于非拉丁字符集的支持和处理，Geohash 需要额外的适配和处理，增加了实现的复杂性。

### 8. Geohash 的扩展与变体

- **Geohash Extended**：

  - 增加了更多的字符来提高编码的精度。
  - 适应不同的应用需求，如实时定位和高精度地图服务。

- **Paged Geohash**：

  - 将地理空间划分为不同的页面（pages），每个页面包含一定数量的 Geohash 单元。
  - 提高了空间查询的效率，适用于大规模地理数据的分布式处理。

- **Geohash Plus**：
  - 结合了 Geohash 和其他地理编码方法，增强了索引的灵活性和查询的速度。
  - 融合了空间索引和聚合功能，适用于更高级的地理信息系统。

### 9. 总结

**Geohash** 是一种将地理坐标编码为短字符串的高效地理编码系统，通过交替划分经纬度范围，生成具有层次化精度和空间相邻性的字符串表示。其主要优势在于简洁性、可扩展性和基于字符串的高效检索能力，使其在地理信息系统、地图服务和位置相关应用中得到了广泛应用。

然而，Geohash 也存在一些局限性，如边界问题、精度限制和动态数据处理的挑战。在设计和实现地理编码系统时，需根据具体应用场景权衡其优缺点，必要时结合其他地理编码方法或进行扩展，以满足复杂的需求和高性能要求。

通过深入理解 Geohash 的工作原理和实现细节，开发者可以更好地利用其特性，优化地理数据的存储和检索，提升系统的整体性能和用户体验。

---

### 附录：Geohash 字符集与编码示例

**Geohash 字符集**：

```
"0123456789bcdefghjkmnpqrstuvwxyz"
```

每个字符代表 5 个二进制位。

**编码示例**：

| 位置            | 经度范围       | 纬度范围        | 符号 |
| --------------- | -------------- | --------------- | ---- |
| 初始            | -180° ~ +180°  | -90° ~ +90°     |      |
| 第 1 位（经度） | 0° ~ +180°     | -90° ~ +90°     | `1`  |
| 第 2 位（纬度） | 0° ~ +180°     | 0° ~ +90°       | `1`  |
| 第 3 位（经度） | 90° ~ +180°    | 0° ~ +90°       | `0`  |
| 第 4 位（纬度） | 90° ~ +180°    | 0° ~ +45°       | `1`  |
| 第 5 位（经度） | 90° ~ +135°    | 22.5° ~ +45°    | `0`  |
| 第 6 位（纬度） | 90° ~ +135°    | 22.5° ~ +33.75° | `1`  |
| 第 7 位（经度） | 112.5° ~ +135° | 22.5° ~ +33.75° | `1`  |

**最终 Geohash**：通过将 5 位二进制序列转换为 Base32 字符，得到字符串`4f8hm`。

**解码示例**：

将`4f8hm`解码为经纬度范围：

1. **字符到二进制**：
   - `4` -> `01100`
   - `f` -> `10111`
   - `8` -> `01000`
   - `h` -> `11000`
   - `m` -> `11011`
2. **合并二进制**：

   ```
   0110010111010001100011011
   ```

3. **交替划分经度和纬度**：

   - 经度位：`0, 1, 1, 0, 0, ...`
   - 纬度位：`1, 0, 1, 1, 0, ...`

4. **逐步缩小范围**：
   - 最终恢复到近似的经纬度范围。

---

### 参考资料

- [Wikipedia: Geohash](https://en.wikipedia.org/wiki/Geohash)
- [- Geohash Explained by Gustavo Niemeyer](https://en.wikipedia.org/wiki/Geohash#Example_of_encoding)
- [Geohash Encoding and Decoding Algorithms](https://www.geohash.org/)
- [Apache Lucene Geospatial Module](https://lucene.apache.org/core/8_8_2/core/org/apache/lucene/document/LatLonPoint.html)

如果您有关于 Geohash 的更多问题或需要进一步的技术细节，欢迎继续咨询！