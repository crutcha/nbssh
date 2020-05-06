package main

import (
	"os"
	"testing"

	"github.com/alecthomas/kingpin"
	"github.com/stretchr/testify/assert"
)

func TestOnlyCommand(t *testing.T) {
	assert := assert.New(t)
	os.Args = []string{"nbssh", "testcommand"}
	kingpin.Parse()

	assert.Equal("testcommand", *command, "command is what we expect")
	assert.Equal("", site.String(), "site is empty")
	assert.Equal("", tenant.String(), "tenant is empty")
	assert.Equal(0, len(*customfield), "custom fields are empty")
}

func TestSingleValueNetboxParam(t *testing.T) {
	site = NetboxParam(kingpin.Flag("site", "Site"))
	assert := assert.New(t)
	os.Args = []string{"nbssh", "testcommand", "--site", "test-site"}
	kingpin.Parse()

	assert.Equal("testcommand", *command, "command is what we expect")
	assert.Equal(1, len(*site), "site contains single entity")
	assert.Equal("test-site", site.String(), "site name is what we expect")
}

func TestMultipleValueNetboxParam(t *testing.T) {
	site = NetboxParam(kingpin.Flag("site", "Site"))
	assert := assert.New(t)
	os.Args = []string{"nbssh", "testcommand", "--site", "test-site-1,test-site-2"}
	kingpin.Parse()

	assert.Equal("testcommand", *command, "command is what we expect")
	assert.Equal(2, len(*site), "site contains multiple entities")
	assert.Equal("test-site-1,test-site-2", site.String(), "site name is what we expect")
	assert.Equal("test-site-1", (*site)[0])
	assert.Equal("test-site-2", (*site)[1])
}

func TestCustomValueParsing(t *testing.T) {
	assert := assert.New(t)
	os.Args = []string{"nbssh", "testcommand", "--customfield", "cf_test=this", "--customfield", "cf_anotha=one"}
	kingpin.Parse()

	assert.Equal("testcommand", *command, "command is what we expect")
	assert.Equal(2, len(*customfield), "custom fields has expected value amount")
	assert.Equal("cf_test", (*customfield)[0].Key)
	assert.Equal("this", (*customfield)[0].Value)
	assert.Equal("cf_anotha", (*customfield)[1].Key)
	assert.Equal("one", (*customfield)[1].Value)
}
