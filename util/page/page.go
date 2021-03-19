package page

const (
	defaultPageSize uint64 = 20
	maxPageSize            = 100
)

func Page(pageSize uint64, pageNo uint64) (uint64, uint64, uint64) {
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageNo == 0 {
		pageNo = 1
	}
	return pageSize, pageNo, (pageNo - 1) * pageSize
}
