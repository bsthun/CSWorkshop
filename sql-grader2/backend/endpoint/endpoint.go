package endpoint

import (
	"backend/common/config"
	"backend/common/fiber/middleware"
	"backend/endpoint/admin"
	"backend/endpoint/public"
	"backend/endpoint/state"
	"backend/endpoint/student"
	"backend/type/common"
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func Bind(
	app *fiber.App,
	config *config.Config,
	frontend common.FrontendFS,
	publicEndpoint *publicEndpoint.Handler,
	stateEndpoint *stateEndpoint.Handler,
	adminEndpoint *adminEndpoint.Handler,
	studentEndpoint *studentEndpoint.Handler,
	middleware *middleware.Middleware,
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
	admin.Post("/collection/detail", adminEndpoint.HandleCollectionDetail)
	admin.Post("/collection/edit", adminEndpoint.HandleCollectionEdit)
	admin.Post("/collection/create", adminEndpoint.HandleCollectionCreate)
	admin.Post("/collection/schema/upload", adminEndpoint.HandleCollectionSchemaUpload)
	admin.Post("/collection/question/create", adminEndpoint.HandleCollectionQuestionCreate)
	admin.Post("/collection/question/edit", adminEndpoint.HandleCollectionQuestionEdit)
	admin.Post("/collection/question/list", adminEndpoint.HandleCollectionQuestionList)
	admin.Post("/collection/question/detail", adminEndpoint.HandleCollectionQuestionDetail)
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
	admin.Post("/exam/edit", adminEndpoint.HandleExamEdit)
	admin.Post("/exam/joinee/list", adminEndpoint.HandleExamJoineeList)
	admin.Post("/exam/question/add", adminEndpoint.HandleExamQuestionAdd)
	admin.Post("/exam/question/list", adminEndpoint.HandleExamQuestionList)
	admin.Post("/exam/question/detail", adminEndpoint.HandleExamQuestionDetail)
	admin.Post("/exam/question/delete", adminEndpoint.HandleExamQuestionDelete)
	admin.Post("/exam/question/edit", adminEndpoint.HandleExamQuestionEdit)
	admin.Post("/submission/detail", adminEndpoint.HandleSubmissionDetail)
	admin.Post("/submission/list", adminEndpoint.HandleSubmissionList)

	// * student endpoints
	student := api.Group("/student", middleware.Jwt(true))
	student.Post("/class/list", studentEndpoint.HandleClassList)
	student.Post("/class/join", studentEndpoint.HandleClassJoin)
	student.Post("/class/exam/list", studentEndpoint.HandleClassExamList)
	student.Post("/class/exam/attempt/detail", studentEndpoint.HandleClassExamAttemptDetail)
	student.Post("/class/exam/attempt", studentEndpoint.HandleClassExamAttempt)
	student.Post("/exam/question/list", studentEndpoint.HandleStudentExamQuestionList)
	student.Post("/exam/question/detail", studentEndpoint.HandleStudentExamQuestionDetail)
	student.Post("/exam/submit", studentEndpoint.HandleExamSubmit)

	// * frontend
	app.Get("*", func(c *fiber.Ctx) error {
		filePath := filepath.Join(".local/dist", c.Path())
		file, err := frontend.ReadFile(filePath)
		if err != nil {
			file, _ = frontend.ReadFile(".local/dist/index.html")
			c.Set("Content-Type", "text/html")
			return c.Send(file)
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		c.Set("Content-Type", contentType)

		return c.Send(file)
	})
}
