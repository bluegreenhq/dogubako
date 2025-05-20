package sql

type Args []any

type Statement interface {
	Query() string
	Args() Args
}

type statement struct {
	query string
	args  Args
}

var _ Statement = (*statement)(nil)

func NewStatement(query string, args Args) *statement {
	return &statement{
		query: query,
		args:  args,
	}
}

func (s *statement) Query() string {
	return s.query
}

func (s *statement) Args() Args {
	return s.args
}
