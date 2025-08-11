package entityProcedure

import (
	"backend/generate/psql"
	"backend/type/payload"
	"context"

	"github.com/bsthun/gut"
)

func (r *Service) ServeSemesterList(ctx context.Context, req *payload.SemesterListRequest) (*payload.SemesterListResponse, *gut.ErrorInstance) {
	// * set default values
	if req.Limit == nil {
		req.Limit = gut.Ptr(int32(10))
	}
	if req.Offset == nil {
		req.Offset = gut.Ptr(int32(0))
	}
	if req.Sort == nil {
		req.Sort = gut.Ptr("createdAt")
	}

	// * count semesters
	count, err := r.database.P().SemesterCount(ctx, req.Name)
	if err != nil {
		return nil, gut.Err(false, "failed to count semesters", err)
	}

	// * list semesters
	semesters, err := r.database.P().SemesterList(ctx, &psql.SemesterListParams{
		Name:   req.Name,
		Sort:   req.Sort,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, gut.Err(false, "failed to list semesters", err)
	}

	// * map to response with classes
	semesterPayloads, er := gut.Iterate(semesters, func(semester psql.Semester) (*payload.Semester, *gut.ErrorInstance) {
		// * get classes for this semester
		classes, err := r.database.P().SemesterClassList(ctx, semester.Id)
		if err != nil {
			return nil, gut.Err(false, "failed to get semester classes", err)
		}

		// * map classes to payload
		classPayloads, _ := gut.Iterate(classes, func(class psql.SemesterClassListRow) (*payload.SemesterClass, *gut.ErrorInstance) {
			return &payload.SemesterClass{
				Id:           class.Id,
				SemesterId:   class.SemesterId,
				Name:         class.Name,
				RegisterCode: class.RegisterCode,
				JoineeCount:  class.JoineeCount,
				CreatedAt:    class.CreatedAt,
				UpdatedAt:    class.UpdatedAt,
			}, nil
		})

		return &payload.Semester{
			Id:        semester.Id,
			Name:      semester.Name,
			Classes:   classPayloads,
			CreatedAt: semester.CreatedAt,
			UpdatedAt: semester.UpdatedAt,
		}, nil
	})
	if er != nil {
		return nil, er
	}

	return &payload.SemesterListResponse{
		Count:     count,
		Semesters: semesterPayloads,
	}, nil
}
