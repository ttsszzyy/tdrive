/*
 * @Author: Young
 * @Date: 2021-05-26 15:07:38
 * @LastEditTime: 2022-09-15 20:39:56
 * @LastEditors: Young
 * @Description: 处理分页查询
 * @FilePath: /buyday/common/utils/tools/page.go
 */

package tools

type Page struct {
	curPage   int64
	size      int64
	totalSzie int64
	totalPage int64
}

func NewPage(curPage int64, size int64, total int64) *Page {
	page := &Page{
		curPage:   curPage,
		size:      size,
		totalSzie: total,
	}

	if page.size > 1000 {
		page.size = 1000
	}

	/*
		page.totalPage = int64(math.Ceil(float64(total) / float64(size)))
		if page.curPage > page.totalPage {
			page.curPage = page.totalPage
		}
		if page.curPage == 0 {
			page.curPage = 1
		}
	*/

	return page
}

func (p *Page) Page() int64 {
	return p.curPage
}

func (p *Page) Offset() int64 {
	return (p.curPage - 1) * p.size
}

func (p *Page) Size() int64 {
	return p.size
}

func (p *Page) TotalSize() int64 {
	return p.totalSzie
}

func (p *Page) TotalPage() int64 {
	return p.totalPage
}
