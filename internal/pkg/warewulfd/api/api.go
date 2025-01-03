package api

import (
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"

	"github.com/warewulf/warewulf/internal/pkg/version"
)

func ApiHandler() *web.Service {
	api := web.NewService(openapi3.NewReflector())

	api.OpenAPISchema().SetTitle("Warewulf v4 API")
	api.OpenAPISchema().SetDescription("This service provides an API to a Warewulf v4 server.")
	api.OpenAPISchema().SetVersion(version.GetVersion())

	api.Get("/api/raw-nodes", apiGetRawNodes())
	api.Get("/api/raw-nodes/{id}", apiGetRawNodeByID())
	api.Put("/api/raw-nodes/{id}", apiPutRawNodeByID())

	api.Get("/api/nodes", apiGetNodes())
	api.Get("/api/nodes/{id}", apiGetNodeByID())

	api.Get("/api/profiles", apiGetProfiles())
	api.Get("/api/profiles/{id}", apiGetProfileByID())

	// container related rest apis
	api.Get("/api/containers", getContainers())
	api.Get("/api/containers/{name}", getContainerByName())
	api.Post("/api/containers/{name}", importContainer())
	api.Delete("/api/containers/{name}", deleteContainer())
	api.Post("/api/containers/{name}/rename/{target}", renameContainer())
	api.Post("/api/containers/{name}/build", buildContainer())

	api.Get("/api/overlays", apiGetOverlays())
	api.Get("/api/overlays/{name}", apiGetOverlayByName())
	api.Get("/api/overlays/{name}/files/{path}", apiGetOverlayFile())

	api.Docs("/api/docs", swgui.New)

	return api
}
