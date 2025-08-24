package payload

import (
	"backend/type/common"
	"backend/type/tuple"
	"time"
)

type CollectionListRequest struct {
	common.Paginate
	Name *string `json:"name"`
}

type Collection struct {
	Id            *uint64                         `json:"id"`
	Name          *string                         `json:"name"`
	Metadata      *tuple.CollectionSchemaMetadata `json:"metadata"`
	CreatedAt     *time.Time                      `json:"createdAt"`
	UpdatedAt     *time.Time                      `json:"updatedAt"`
	QuestionCount *uint64                         `json:"questionCount"`
}

type CollectionListResponse struct {
	Count       *uint64       `json:"count"`
	Collections []*Collection `json:"collections"`
}

type SemesterListRequest struct {
}

type SemesterClass struct {
	Id           *uint64    `json:"id"`
	SemesterId   *uint64    `json:"semesterId"`
	Name         *string    `json:"name"`
	RegisterCode *string    `json:"registerCode"`
	JoineeCount  *uint64    `json:"joineeCount"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type Semester struct {
	Id        *uint64          `json:"id"`
	Name      *string          `json:"name"`
	Classes   []*SemesterClass `json:"classes"`
	CreatedAt *time.Time       `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt"`
}

type SemesterListResponse struct {
	Count     *uint64     `json:"count"`
	Semesters []*Semester `json:"semesters"`
}

type SemesterCreateRequest struct {
	Name *string `json:"name" validate:"required"`
}

type SemesterEditRequest struct {
	Id   *uint64 `json:"id" validate:"required"`
	Name *string `json:"name" validate:"required"`
}

type ClassCreateRequest struct {
	SemesterId *uint64 `json:"semesterId" validate:"required"`
	Code       *string `json:"code" validate:"required"`
	Name       *string `json:"name" validate:"required"`
}

type ClassEditRequest struct {
	Id           *uint64 `json:"id" validate:"required"`
	Code         *string `json:"code" validate:"required"`
	Name         *string `json:"name" validate:"required"`
	RegisterCode *string `json:"registerCode" validate:"required"`
}

type CollectionCreateRequest struct {
	Name *string `json:"name" validate:"required"`
}

type CollectionEditRequest struct {
	Id   *uint64 `json:"id" validate:"required"`
	Name *string `json:"name" validate:"required"`
}

type CollectionSchemaUploadRequest struct {
	CollectionId string `json:"collectionId" validate:"required"`
}

type CollectionQuestionCreateRequest struct {
	CollectionId *uint64 `json:"collectionId" validate:"required"`
	Description  *string `json:"description"`
}

type CollectionQuestionEditRequest struct {
	Id          *uint64 `json:"id" validate:"required"`
	Title       *string `json:"title" validate:"required"`
	Description *string `json:"description"`
	CheckQuery  *string `json:"checkQuery" validate:"required"`
	CheckPrompt *string `json:"checkPrompt" validate:"required"`
}

type CollectionQuestionListResponse struct {
	Questions []*CollectionQuestionItem `json:"questions"`
}

type CollectionQuestionItem struct {
	Id          *uint64 `json:"id"`
	OrderNum    *int32  `json:"orderNum"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CollectionQuestionDetail struct {
	Id           *uint64    `json:"id"`
	CollectionId *uint64    `json:"collectionId"`
	OrderNum     *int32     `json:"orderNum"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	CheckQuery   *string    `json:"checkQuery"`
	CheckPrompt  *string    `json:"checkPrompt"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type CollectionQuestionDeleteRequest struct {
	Id *uint64 `json:"id" validate:"required"`
}

type Class struct {
	Id           *uint64    `json:"id"`
	SemesterId   *uint64    `json:"semesterId"`
	Code         *string    `json:"code"`
	Name         *string    `json:"name"`
	RegisterCode *string    `json:"registerCode"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type User struct {
	Id         *uint64    `json:"id"`
	Oid        *string    `json:"oid"`
	Firstname  *string    `json:"firstname"`
	Lastname   *string    `json:"lastname"`
	Email      *string    `json:"email"`
	PictureUrl *string    `json:"pictureUrl"`
	IsAdmin    *bool      `json:"isAdmin"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

type ClassJoinee struct {
	User      *User      `json:"user"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type SemesterInfo struct {
	Id        *uint64    `json:"id"`
	Name      *string    `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type ClassDetailResponse struct {
	Class    *Class         `json:"class"`
	Semester *SemesterInfo  `json:"semester"`
	Joinees  []*ClassJoinee `json:"joinees"`
}

type ClassJoineeListResponse struct {
	Joinees []*ClassJoinee `json:"joinees"`
}

type ExamCreateRequest struct {
	ClassId      *uint64    `json:"classId" validate:"required"`
	CollectionId *uint64    `json:"collectionId" validate:"required"`
	Name         *string    `json:"name" validate:"required"`
	OpenedAt     *time.Time `json:"openedAt" validate:"required"`
	ClosedAt     *time.Time `json:"closedAt" validate:"required"`
}

type Exam struct {
	Id            *uint64    `json:"id"`
	ClassId       *uint64    `json:"classId"`
	CollectionId  *uint64    `json:"collectionId"`
	Name          *string    `json:"name"`
	AccessCode    *string    `json:"accessCode"`
	OpenedAt      *time.Time `json:"openedAt"`
	ClosedAt      *time.Time `json:"closedAt"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	QuestionCount *uint64    `json:"questionCount"`
}

type ExamListItem struct {
	Exam       *Exam       `json:"exam"`
	Collection *Collection `json:"collection"`
}

type ExamListResponse struct {
	Exams []*ExamListItem `json:"exams"`
}

type ExamJoinee struct {
	Id   *uint64 `json:"id"`
	User *User   `json:"user"`
}

type ExamScore struct {
	Passed      *uint64 `json:"passed"`
	Rejected    *uint64 `json:"rejected"`
	Invalid     *uint64 `json:"invalid"`
	Unsubmitted *int32  `json:"unsubmitted"`
}

type ExamJoineeListItem struct {
	Id            *uint64     `json:"id"`
	ExamId        *uint64     `json:"examId"`
	ClassJoineeId *uint64     `json:"classJoineeId"`
	OpenedAt      *time.Time  `json:"openedAt"`
	StartedAt     *time.Time  `json:"startedAt"`
	FinishedAt    *time.Time  `json:"finishedAt"`
	CreatedAt     *time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time  `json:"updatedAt"`
	Joinee        *ExamJoinee `json:"joinee"`
	Score         *ExamScore  `json:"score"`
}

type ExamJoineeListResponse struct {
	Joinees []*ExamJoineeListItem `json:"joinees"`
}

type ExamQuestion struct {
	Id                 *uint64    `json:"id"`
	ExamId             *uint64    `json:"examId"`
	OriginalQuestionId *uint64    `json:"originalQuestionId"`
	OrderNum           *int32     `json:"orderNum"`
	Title              *string    `json:"title"`
	Description        *string    `json:"description"`
	CheckQuery         *string    `json:"checkQuery"`
	CheckPrompt        *string    `json:"checkPrompt"`
	CreatedAt          *time.Time `json:"createdAt"`
	UpdatedAt          *time.Time `json:"updatedAt"`
}

type ExamAttempt struct {
	Id            *uint64    `json:"id"`
	ExamId        *uint64    `json:"examId"`
	ClassJoineeId *uint64    `json:"classJoineeId"`
	OpenedAt      *time.Time `json:"openedAt"`
	StartedAt     *time.Time `json:"startedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

type ExamSubmission struct {
	Id                *uint64    `json:"id"`
	ExamQuestionId    *uint64    `json:"examQuestionId"`
	ExamAttemptId     *uint64    `json:"examAttemptId"`
	Answer            *string    `json:"answer"`
	CheckQueryPassed  *bool      `json:"checkQueryPassed"`
	CheckQueryAt      *time.Time `json:"checkQueryAt"`
	CheckPromptPassed *bool      `json:"checkPromptPassed"`
	CheckPromptAt     *time.Time `json:"checkPromptAt"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}

type SubmissionDetailResponse struct {
	Submission *ExamSubmission `json:"submission"`
	Question   *ExamQuestion   `json:"question"`
	Attempt    *ExamAttempt    `json:"attempt"`
}

type ClassJoineeInfo struct {
	Id        *uint64    `json:"id"`
	UserId    *uint64    `json:"userId"`
	ClassId   *uint64    `json:"classId"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type SubmissionListItem struct {
	Submission *ExamSubmission  `json:"submission"`
	Question   *ExamQuestion    `json:"question"`
	Attempt    *ExamAttempt     `json:"attempt"`
	Student    *User            `json:"student"`
	Joinee     *ClassJoineeInfo `json:"joinee"`
}

type SubmissionIdRequest struct {
	SubmissionId *uint64 `json:"submissionId" validate:"required"`
}

type SubmissionListRequest struct {
	ExamAttemptId  *uint64 `json:"examAttemptId"`
	ExamQuestionId *uint64 `json:"examQuestionId"`
}

type SubmissionListResponse struct {
	Exam        *Exam                 `json:"exam"`
	Submissions []*SubmissionListItem `json:"submissions"`
}

type ExamQuestionAddRequest struct {
	ExamId               *uint64 `json:"examId" validate:"required"`
	CollectionQuestionId *uint64 `json:"collectionQuestionId" validate:"required"`
}

type ExamQuestionDeleteRequest struct {
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
}

type ExamQuestionEditRequest struct {
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
	Title          *string `json:"title" validate:"required"`
	Description    *string `json:"description"`
	CheckQuery     *string `json:"checkQuery" validate:"required"`
	CheckPrompt    *string `json:"checkPrompt" validate:"required"`
}

type ExamQuestionListResponse struct {
	Questions []*CollectionQuestionItem `json:"questions"`
}

type ClassIdRequest struct {
	ClassId *uint64 `json:"classId" validate:"required"`
}

type ExamIdRequest struct {
	ExamId *uint64 `json:"examId" validate:"required"`
}

type CollectionIdRequest struct {
	CollectionId *uint64 `json:"collectionId" validate:"required"`
}

type QuestionIdRequest struct {
	QuestionId *uint64 `json:"questionId" validate:"required"`
}

type CollectionDetailResponse struct {
	Collection *Collection `json:"collection"`
}

type ExamAttemptCount struct {
	OpenedCount   *uint64 `json:"openedCount"`
	StartedCount  *uint64 `json:"startedCount"`
	FinishedCount *uint64 `json:"finishedCount"`
}

type ExamQuestionIdRequest struct {
	ExamQuestionId *uint64 `json:"examQuestionId" validate:"required"`
}

type ExamQuestionDetailResponse struct {
	ExamQuestion       *ExamQuestion             `json:"examQuestion"`
	CollectionQuestion *CollectionQuestionDetail `json:"collectionQuestion"`
}

type ExamDetailResponse struct {
	Exam         *Exam             `json:"exam"`
	Class        *Class            `json:"class"`
	Collection   *Collection       `json:"collection"`
	AttemptCount *ExamAttemptCount `json:"attemptCount"`
}

type StudentClassListItem struct {
	Class             *Class        `json:"class"`
	Semester          *SemesterInfo `json:"semester"`
	ExamTotalCount    *uint64       `json:"examTotalCount"`
	ExamFinishedCount *uint64       `json:"examFinishedCount"`
}

type StudentClassListRequest struct {
}

type StudentClassListResponse struct {
	Classes []*StudentClassListItem `json:"classes"`
}

type ClassJoinRequest struct {
	RegisterCode *string `json:"registerCode" validate:"required"`
}

type ClassExamListRequest struct {
	ClassId *uint64 `json:"classId" validate:"required"`
}

type ClassExamListItem struct {
	Exam          *Exam   `json:"exam"`
	QuestionCount *uint64 `json:"questionCount"`
}

type ClassExamListResponse struct {
	Exams []*ClassExamListItem `json:"exams"`
}

type ClassExamAttemptRequest struct {
	ExamId     *uint64 `json:"examId" validate:"required"`
	AccessCode *string `json:"accessCode" validate:"required"`
}

type ClassExamDetailRequest struct {
	ExamId *uint64 `json:"examId" validate:"required"`
}

type ExamCredential struct {
	Dialect      *string `json:"dialect"`
	Host         *string `json:"host"`
	Port         *string `json:"port"`
	User         *string `json:"user"`
	Password     *string `json:"password"`
	DatabaseName *string `json:"databaseName"`
}

type ClassExamDetailResponse struct {
	Class             *Class          `json:"class"`
	Exam              *Exam           `json:"exam"`
	ExamQuestionCount *uint64         `json:"examQuestionCount"`
	Credential        *ExamCredential `json:"credential"`
}
