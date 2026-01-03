package dbrepo

import (
	"fmt"
	"slices"
	"strings"
)

var (
	// 限制最小页码
	minPageNum = 1

	// 默认每页大小
	defaultPageSize = 10

	// 限制每页最大条数
	maxPageSize = 100
)

type Filters struct {
	PageNum      int
	PageSize     int
	SortFields   []string // SQL ORDER BY 排序字段, 必须存在于 SortSafelist
	SortSafelist []string // 定义安全的排序字段，带"-"的是DESC，反之ASC
}

// 使用时，先要设置 SortSafelist 安全字段值
//
// sortColumn 获取安全的排序字段，有"-"就去掉,返回一个string数组：
// ['created_at DESC', 'id ASC']
func (f Filters) sortColumn(br baseRepo) string {
	if len(f.SortSafelist) == 0 {
		f.SortSafelist = br.defaultSortSafelist()
	}

	sorts := make([]string, 0, len(f.SortFields))

	for _, field := range f.SortFields {
		if slices.Contains(f.SortSafelist, field) {
			name := strings.TrimPrefix(field, "-")
			sort := f.sortDirection(field)
			sorts = append(sorts, fmt.Sprintf("%s %s", name, sort))
		}
	}

	return strings.Join(sorts, ",")
}

// sortColumnWithDefault 默认值为 "id ASC" 的 sortColumn 方法
func (f Filters) sortColumnWithDefault(br baseRepo) string {
	sortText := f.sortColumn(br)
	if sortText == "" {
		return " id ASC "
	}

	return sortText
}

// sortDirection 获取排序顺序 带"-"为DESC，反之ASC
func (f Filters) sortDirection(sortField string) string {
	if strings.HasPrefix(sortField, "-") {
		return "DESC"
	}

	return "ASC"
}

// check 检查PageSize、PageNum 是否满足条件，不满足就设置为默认值
func (f *Filters) check() {
	if f.PageSize <= 0 {
		f.PageSize = defaultPageSize
	}

	if f.PageSize > maxPageSize {
		f.PageSize = maxPageSize
	}

	if f.PageNum <= 0 {
		f.PageNum = minPageNum
	}
}

// limit 先检查在返回限制每页多少条
func (f Filters) limit() int {
	f.check()
	return f.PageSize
}

// offset 先检查在返回分页起点
func (f Filters) offset() int {
	f.check()
	return (f.PageNum - 1) * f.PageSize
}

type Metadata struct {
	PageNum    int `json:"pageNum,omitzero"`
	PageSize   int `json:"pageSize,omitzero"`
	LastPage   int `json:"lastPage,omitzero"`
	TotalCount int `json:"totalCount,omitzero"`
}

// PageQueryVo 分页数据通用结构体，意在分页数据返回统一结构体
type PageQueryVo struct {
	List      any      `json:"list"`
	Metadata  Metadata `json:"metadata"`
	ExtraData any      `json:"extraData,omitempty"`
}
