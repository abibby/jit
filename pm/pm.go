package pm

type Provider interface {
	SetStatus(issueID, status string) error
	GetIssue(issueID string) (*Issue, error)
	GetMyIssues() ([]*Issue, error)
}

type Status string

var (
	StatusToDo       = Status("to-do")
	StatusInProgress = Status("in-progress")
	StatusInReview   = Status("in-review")
	StatusDone       = Status("done")
	StatusUnknown    = Status("unknown")
)

type Issue struct {
	ID     string
	Title  string
	Status Status
}

func GetProvider() (Provider, error) {
	return NewJira(), nil
}
