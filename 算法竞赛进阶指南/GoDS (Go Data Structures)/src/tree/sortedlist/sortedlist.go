// !fhq-treap

package sortedlist

type SortedList struct {
	comparator Comparator
}

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b interface{}) int

func NewSortedList(comparator Comparator) *SortedList {}

func (sl *SortedList) Add(value interface{})     {}
func (sl *SortedList) Has(value interface{})     {}
func (sl *SortedList) At(index int) interface{}  {}
func (sl *SortedList) Pop(index int) interface{} {}
func (sl *SortedList) Remove(value interface{})  {}
func (sl *SortedList) Discard(value interface{}) {}

// 查询排名
func (sl *SortedList) BisectLeft(value interface{}) int  {}
func (sl *SortedList) BisectRight(value interface{}) int {}

func (sl *SortedList) Len() int {}

func (sl *SortedList) String() string {}

func (sl *SortedList) lower(value interface{}) interface{} {}
func (sl *SortedList) upper(value interface{}) interface{} {}
