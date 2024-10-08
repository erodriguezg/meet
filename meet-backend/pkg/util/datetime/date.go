package datetime

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

const dateLayout = "2006-01-02"

// UnmarshalJSON Parses the json string in the custom format
func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(dateLayout, s)
	*ct = Date(nt)
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct Date) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *Date) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(dateLayout))
}
