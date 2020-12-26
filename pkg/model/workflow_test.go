package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWorkflow_StringEvent(t *testing.T) {
	yaml := `
name: local-action-docker-url
on: push

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")

	assert.Len(t, workflow.On(), 1)
	assert.Contains(t, workflow.On(), "push")
	assert.Empty(t, workflow.OnPaths("push"))
}

func TestReadWorkflow_ListEvent(t *testing.T) {
	yaml := `
name: local-action-docker-url
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")

	assert.Len(t, workflow.On(), 2)
	assert.Contains(t, workflow.On(), "push")
	assert.Contains(t, workflow.On(), "pull_request")
	assert.Empty(t, workflow.OnPaths("push"))
	assert.Empty(t, workflow.OnPaths("pull_request"))
}

func TestReadWorkflow_MapEvent(t *testing.T) {
	yaml := `
name: local-action-docker-url
on: 
  push:
    branches:
    - master
  pull_request:
    branches:
    - master

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")
	assert.Len(t, workflow.On(), 2)
	assert.Contains(t, workflow.On(), "push")
	assert.Contains(t, workflow.On(), "pull_request")
	assert.Empty(t, workflow.OnPaths("push"))
	assert.Empty(t, workflow.OnPaths("pull_request"))
}

func TestReadWorkflow_MapEventPaths(t *testing.T) {
	yaml := `
name: local-action-docker-url
on: 
  push:
    branches:
    - master
    paths:
    - "a*"
    - "b*"
    - "!c*"
  pull_request:
    branches:
    - master

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")
	assert.Len(t, workflow.On(), 2)
	assert.Contains(t, workflow.On(), "push")
	assert.Contains(t, workflow.On(), "pull_request")
	assert.Contains(t, workflow.OnPaths("push"), "a*")
	assert.Contains(t, workflow.OnPaths("push"), "b*")
	assert.Contains(t, workflow.OnPaths("push"), "!c*")
	assert.Empty(t, workflow.OnPaths("pull_request"))
}

func TestReadWorkflow_MapEventPathsIgnore(t *testing.T) {
	yaml := `
name: local-action-docker-url
on: 
  push:
    branches:
    - master
    paths-ignore:
    - "a*"
    - "b*"
    - "c*"
  pull_request:
    branches:
    - master

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")
	assert.Len(t, workflow.On(), 2)
	assert.Contains(t, workflow.On(), "push")
	assert.Contains(t, workflow.On(), "pull_request")
	assert.Contains(t, workflow.OnPaths("push"), "!a*")
	assert.Contains(t, workflow.OnPaths("push"), "!b*")
	assert.Contains(t, workflow.OnPaths("push"), "!c*")
	assert.Empty(t, workflow.OnPaths("pull_request"))
}
func TestReadWorkflow_StringContainer(t *testing.T) {
	yaml := `
name: local-action-docker-url

jobs:
  test:
    container: nginx:latest
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
  test2:
    container: 
      image: nginx:latest
      env:
        foo: bar
    runs-on: ubuntu-latest
    steps:
    - uses: ./actions/docker-url
`

	workflow, err := ReadWorkflow(strings.NewReader(yaml))
	assert.NoError(t, err, "read workflow should succeed")
	assert.Len(t, workflow.Jobs, 2)
	assert.Contains(t, workflow.Jobs["test"].Container().Image, "nginx:latest")
	assert.Contains(t, workflow.Jobs["test2"].Container().Image, "nginx:latest")
	assert.Contains(t, workflow.Jobs["test2"].Container().Env["foo"], "bar")
}
