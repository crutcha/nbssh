package main

import (
    "os"
	"github.com/stretchr/testify/assert"
	"testing"
    "github.com/alecthomas/kingpin"
)

func TestOnlyCommand(t *testing.T) {
    assert := assert.New(t)
    os.Args = []string{"nbssh", "testcommand"}
    kingpin.Parse()

    assert.Equal(*command, "testcommand", "command is what we expect")
    assert.Equal(site.String(), "", "site is empty")
    assert.Equal(tenant.String(), "", "tenant is empty")
    assert.Equal(len(*customfield), 0, "custom fields are empty")
}
