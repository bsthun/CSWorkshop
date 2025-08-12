package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"strconv"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get class id from params
	classIdParam := c.Params("id")
	classId, err := strconv.ParseUint(classIdParam, 10, 64)
	if err != nil {
		return gut.Err(false, "invalid class id", err)
	}

	// * get user from database
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user is admin
	if user.IsAdmin == nil || !*user.IsAdmin {
		return gut.Err(false, "access denied", nil)
	}

	// * get class with semester info
	classDetail, err := r.database.P().ClassGetById(c.Context(), &classId)
	if err != nil {
		return gut.Err(false, "failed to get class", err)
	}

	// * list class joinees
	joineeRows, err := r.database.P().ClassJoineeList(c.Context(), &classId)
	if err != nil {
		return gut.Err(false, "failed to get class joinees", err)
	}

	// * map joinees to payload
	joinees, er := gut.Iterate(joineeRows, func(row psql.ClassJoineeListRow) (*payload.ClassJoinee, *gut.ErrorInstance) {
		return &payload.ClassJoinee{
			User: &payload.User{
				Id:         row.User.Id,
				Oid:        row.User.Oid,
				Firstname:  row.User.Firstname,
				Lastname:   row.User.Lastname,
				Email:      row.User.Email,
				PictureUrl: row.User.PictureUrl,
				IsAdmin:    row.User.IsAdmin,
				CreatedAt:  row.User.CreatedAt,
				UpdatedAt:  row.User.UpdatedAt,
			},
			CreatedAt: row.ClassJoinee.CreatedAt,
			UpdatedAt: row.ClassJoinee.UpdatedAt,
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.ClassDetailResponse{
		Class: &payload.Class{
			Id:           classDetail.Class.Id,
			SemesterId:   classDetail.Class.SemesterId,
			Code:         classDetail.Class.Code,
			Name:         classDetail.Class.Name,
			RegisterCode: classDetail.Class.RegisterCode,
			CreatedAt:    classDetail.Class.CreatedAt,
			UpdatedAt:    classDetail.Class.UpdatedAt,
		},
		Semester: &payload.SemesterInfo{
			Id:        classDetail.Semester.Id,
			Name:      classDetail.Semester.Name,
			CreatedAt: classDetail.Semester.CreatedAt,
			UpdatedAt: classDetail.Semester.UpdatedAt,
		},
		Joinees: joinees,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
