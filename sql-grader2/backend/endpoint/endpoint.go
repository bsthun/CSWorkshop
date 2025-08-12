package endpoint

import (
	"backend/common/config"
	"backend/common/fiber/middleware"
	"backend/endpoint/admin"
	"backend/endpoint/public"
	"backend/endpoint/state"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func Bind(
	app *fiber.App,
	publicEndpoint *publicEndpoint.Handler,
	stateEndpoint *stateEndpoint.Handler,
	adminEndpoint *adminEndpoint.Handler,
	middleware *middleware.Middleware,
	config *config.Config,
) {
	api := app.Group("/api")
	api.Use(middleware.Id())

	// * public endpoints
	public := api.Group("/public")
	public.Get("/login/redirect", publicEndpoint.HandleLoginRedirect)
	public.Post("/login/callback", publicEndpoint.HandleLoginCallback)

	// * state endpoints
	state := api.Group("/state", middleware.Jwt(true))
	state.Post("/state", stateEndpoint.HandleState)

	// * admin endpoints
	admin := api.Group("/admin", middleware.Jwt(true))
	admin.Post("/collection/list", adminEndpoint.HandleCollectionList)
	admin.Post("/collection/create", adminEndpoint.HandleCollectionCreate)
	admin.Post("/collection/schema/upload", adminEndpoint.HandleCollectionSchemaUpload)
	admin.Post("/collection/question/create", adminEndpoint.HandleCollectionQuestionCreate)
	admin.Post("/collection/question/edit", adminEndpoint.HandleCollectionQuestionEdit)
	admin.Get("/collection/question/list", adminEndpoint.HandleCollectionQuestionList)
	admin.Get("/collection/question/detail", adminEndpoint.HandleCollectionQuestionDetail)
	admin.Post("/collection/question/delete", adminEndpoint.HandleCollectionQuestionDelete)
	admin.Post("/semester/list", adminEndpoint.HandleSemesterList)
	admin.Post("/semester/create", adminEndpoint.HandleSemesterCreate)
	admin.Post("/semester/edit", adminEndpoint.HandleSemesterEdit)
	admin.Post("/class/create", adminEndpoint.HandleClassCreate)
	admin.Post("/class/detail", adminEndpoint.HandleClassDetail)
	admin.Post("/class/edit", adminEndpoint.HandleClassEdit)
	admin.Post("/exam/create", adminEndpoint.HandleExamCreate)
	admin.Post("/exam/list", adminEndpoint.HandleExamList)
	admin.Post("/exam/detail", adminEndpoint.HandleExamDetail)
	admin.Post("/exam/joinee/list", adminEndpoint.HandleExamJoineeList)
	admin.Post("/exam/question/add", adminEndpoint.HandleExamQuestionAdd)
	admin.Post("/exam/question/delete", adminEndpoint.HandleExamQuestionDelete)
	admin.Post("/exam/question/edit", adminEndpoint.HandleExamQuestionEdit)
	admin.Post("/exam/question/list", adminEndpoint.HandleExamQuestionList)
	admin.Post("/submission/detail", adminEndpoint.HandleSubmissionDetail)
	admin.Post("/submission/list", adminEndpoint.HandleSubmissionList)

	// * static files
	app.Static("/file", ".local/file")

	// * static
	app.Static("/", *config.WebRoot)
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(*config.WebRoot, "index.html"))
	})
}
