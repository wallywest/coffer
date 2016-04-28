package recording

import (
	"fmt"

	"github.com/nats-io/nuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CURSOR_PAGE_SIZE = 20
)

type cursorQuery struct {
	Session    *mgo.Session
	DB         string
	Collection string
	Selector   bson.M
}

type mongoCursor struct {
	Id       string
	Selector bson.M
	Query    *cursorQuery

	iter       *mgo.Iter
	collection *mgo.Collection

	pageSize int
	total    int
	current  int
	hasNext  bool
}

type CursorInfo struct {
	CursorId  string
	Total     int
	Start     int
	End       int
	PageSize  int
	Numpages  int
	PageIndex int
}

func newRecordingCursor(query *cursorQuery) (*mongoCursor, error) {
	nuid := nuid.Next()

	collection := query.Session.DB(query.DB).C(query.Collection)
	q := collection.Find(query.Selector).Batch(CURSOR_PAGE_SIZE)
	count, err := q.Count()

	if err != nil {
		return nil, err
	}

	return &mongoCursor{
		Id:       nuid,
		Query:    query,
		pageSize: CURSOR_PAGE_SIZE,
		total:    count,
		iter:     q.Iter(),
		current:  0,
		hasNext:  true,
	}, nil
}

func (c *mongoCursor) PageSize() int {
	return c.pageSize
}

func (c *mongoCursor) SetPageSize(size int) {
}

func (c *mongoCursor) Current() int {
	return c.current
}

func (c *mongoCursor) HasNext() bool {
	return c.hasNext
}

func (c *mongoCursor) Total() int {
	return c.total
}

func (c *mongoCursor) Close() error {
	err := c.iter.Close()
	c.Query.Session.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *mongoCursor) NextPage() ([]*Recording, *CursorInfo, error) {
	var recordings []*Recording

	start := c.current

	for i := 0; i < c.pageSize; i++ {

		if c.total == 0 {
			c.hasNext = false
			break
		}

		if !c.hasNext {
			break
		}

		var rec Recording

		iterResult := c.iter.Next(&rec)

		err := c.iter.Err()

		if err != nil {
			c.hasNext = false
			return nil, nil, err
		}

		timeout := c.iter.Timeout()

		if timeout == true {
			c.hasNext = false
			return nil, nil, fmt.Errorf("timeout on on cursor")
		}

		if iterResult == false {
			c.hasNext = false
			break
		}

		recordings = append(recordings, &rec)
		c.current += 1
	}

	var end int
	if c.current > start {
		end = c.current - 1
	} else {
		end = start
	}

	p := &CursorInfo{
		CursorId:  c.Id,
		Total:     c.total,
		Start:     start,
		End:       end,
		Numpages:  (c.total + c.pageSize - 1) / c.pageSize,
		PageIndex: start / c.pageSize,
		PageSize:  c.pageSize,
	}

	return recordings, p, nil
}
