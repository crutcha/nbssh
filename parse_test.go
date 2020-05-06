package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
