package mesos_data

import (
	"fmt"
	"net"
	"strings"
	"github.com/emc-advanced-dev/pkg/errors"
)

// UPID is a equivalent of the UPID in libprocess.
type UPID struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// Parse parses the UPID from the input string.
func UPIDFromString(input string) (*UPID, error) {
	upid := new(UPID)

	splits := strings.Split(input, "@")
	if len(splits) != 2 {
		return nil, fmt.Errorf("Expect one `@' in the input")
	}
	upid.ID = splits[0]

	var err error
	upid.Host, upid.Port, err = net.SplitHostPort(splits[1])
	if err != nil {
		return nil, errors.New("failed to split host and port", err)
	}
	return upid, nil
}

// String returns the string representation.
func (u UPID) String() string {
	return fmt.Sprintf("%s@%s:%s", u.ID, u.Host, u.Port)
}

// Equal returns true if two upid is equal
func (u *UPID) Equal(upid *UPID) bool {
	if u == nil {
		return upid == nil
	} else {
		return upid != nil && u.ID == upid.ID && u.Host == upid.Host && u.Port == upid.Port
	}
}
