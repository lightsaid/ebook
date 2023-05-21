package dbrepo

type Pager struct {
	size int // 每页多条
	page int // 当前第几页
}

// NewPager 创建一个分页器 page >= 1, <=10 size <= 100
func NewPager(page int, size int) Pager {
	if size > 100 {
		size = 100
	}
	if size < 10 {
		size = 10
	}
	if page < 1 {
		page = 1
	}
	return Pager{
		size: size,
		page: page,
	}
}

func (p *Pager) Limit() int {
	return p.size
}
func (p *Pager) Offset() int {
	return (p.page - 1) * p.size
}
