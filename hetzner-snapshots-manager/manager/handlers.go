package manager

import (
	"errors"
	"net/http"

	"github.com/spigell/pulumi-automation-api-apps/common/apiserver"
	"github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager/hetzner"

	"github.com/gin-gonic/gin"
)

const (
	SnapshotsAPIPath = "/hetzner/snapshots"
	ErrServerMissing = "`server` parameter required"
)

// Additional handlers for api server.
func getAllRoutes(snapshots *hetzner.Snapshots) []apiserver.Route {
	return []apiserver.Route{
		{
			Path: SnapshotsAPIPath,
			Handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					server := c.Query("server")
					if server == "" {
						c.JSON(http.StatusNotFound, gin.H{
							"error":  ErrServerMissing,
							"status": "ERROR",
						})
					}
					lastSnapshotInfo, err := snapshots.GetLast(server)
					if err != nil {
						if errors.Is(err, hetzner.ErrSnapshotNotFound) {
							c.JSON(http.StatusNotFound, gin.H{
								"error":  err.Error(),
								"status": "ERROR",
							})

							return
						}
						c.JSON(http.StatusInternalServerError, gin.H{
							"error":  err.Error(),
							"status": "ERROR",
						})

						return
					}

					c.JSON(http.StatusOK, gin.H{
						"body":   lastSnapshotInfo,
						"status": "OK",
					})
				}
			}(),
		},
	}
}
