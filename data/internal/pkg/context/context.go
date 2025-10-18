package context

import (
	"context"
	"fmt"

	"github.com/go-mockingcode/data/internal/pkg/project"
)

type contextKey string

const ProjectKey contextKey = "project"

func GetProjectInfo(ctx context.Context) (*project.ProjectInfo, error) {
	projectInfo, ok := ctx.Value(ProjectKey).(*project.ProjectInfo)
	if !ok {
		return nil, fmt.Errorf("project not found in context")
	}
	return projectInfo, nil
}
