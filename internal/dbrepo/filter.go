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

// sortColumn 获取安全的排序字段，有"-"就去掉,返回一个string数组：
// ['created_at DESC', 'id ASC']
func (f Filters) sortColumn() string {
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

// limit 限制每页多少条
func (f Filters) limit() int {
	f.check()
	return f.PageSize
}

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

func calculateMetadata(totalCount, pageNum, pageSize int) Metadata {
	if totalCount == 0 {
		return Metadata{}
	}

	return Metadata{
		PageNum:    pageNum,
		PageSize:   pageSize,
		LastPage:   (totalCount + pageSize - 1) / pageSize,
		TotalCount: totalCount,
	}
}

type PageQueryVo struct {
	List     any      `json:"list"`
	Metadata Metadata `json:"metadata"`
}

// 构建统一返回数据
func makePageQueryVo(metadata Metadata, list any) *PageQueryVo {
	return &PageQueryVo{
		List:     list,
		Metadata: metadata,
	}
}
