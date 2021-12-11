package papermc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL = "https://papermc.io/api/v2"

type Version struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

func GetVersion(ctx context.Context, project, version string) (Version, error) {
	route := fmt.Sprintf("/projects/%s/versions/%s", project, version)
	r, err := http.Get(URL + route)
	if err != nil {
		return Version{}, err
	}
	defer r.Body.Close()
	var v Version
	err = json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func Download(ctx context.Context, project, version string, build int) (io.ReadCloser, error) {
	download := fmt.Sprintf("%s-%s-%d.jar", project, version, build)
	route := fmt.Sprintf("/projects/%s/versions/%s/builds/%d/downloads/%s", project, version, build, download)
	r, err := http.Get(URL + route)
	if err != nil {
		return nil, err
	}
	return r.Body, err
}
