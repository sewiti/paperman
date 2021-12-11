package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAll(t *testing.T) {
	expected := []Server{
		{
			Name:     "gamma",
			Version:  "1.18",
			Port:     25565,
			JavaArgs: []string{"-Xms1500M", "-Xmx1500M"},
			JarArgs:  []string{"nogui"},
		},
		{
			Name:     "johhny",
			Version:  "1.18.1",
			Port:     25588,
			JavaArgs: []string{"-Xms1500M", "-Xmx1500M"},
			JarArgs:  []string{"nogui"},
		},
	}

	got, err := ReadAll("testdata")
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
