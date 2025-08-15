-- +goose Up
-- +goose StatementBegin

-- * reuse users table from original schema
CREATE TABLE users
(
    id          BIGSERIAL PRIMARY KEY,
    oid         VARCHAR(64)  NOT NULL UNIQUE,
    firstname   VARCHAR(255) NOT NULL,
    lastname    VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE,
    picture_url TEXT         NULL,
    is_admin    BOOLEAN      NOT NULL,
    metadata    JSONB        NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE collections
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    metadata   JSONB        NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE collection_questions
(
    id            BIGSERIAL PRIMARY KEY,
    collection_id BIGINT REFERENCES collections (id) ON DELETE CASCADE NOT NULL,
    order_num     INTEGER                                              NOT NULL,
    title         VARCHAR(255)                                         NOT NULL,
    description   TEXT                                                 NULL,
    check_query   TEXT                                                 NOT NULL,
    check_prompt  TEXT                                                 NOT NULL,
    created_at    TIMESTAMP                                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP                                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (collection_id, order_num)
);

CREATE TABLE semesters
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE classes
(
    id            BIGSERIAL PRIMARY KEY,
    semester_id   BIGINT REFERENCES semesters (id) ON DELETE CASCADE NOT NULL,
    code          VARCHAR(255)                                       NOT NULL,
    name          VARCHAR(255)                                       NOT NULL,
    register_code VARCHAR(64)                                        NOT NULL UNIQUE,
    created_at    TIMESTAMP                                          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP                                          NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE class_joinees
(
    id         BIGSERIAL PRIMARY KEY,
    class_id   BIGINT REFERENCES classes (id) ON DELETE CASCADE NOT NULL,
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE   NOT NULL,
    created_at TIMESTAMP                                        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                                        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (class_id, user_id)
);

CREATE TABLE exams
(
    id            BIGSERIAL PRIMARY KEY,
    class_id      BIGINT REFERENCES classes (id) ON DELETE CASCADE      NOT NULL,
    collection_id BIGINT REFERENCES collections (id) ON DELETE RESTRICT NOT NULL,
    name          VARCHAR(255)                                          NOT NULL,
    access_code   VARCHAR(64)                                           NOT NULL,
    opened_at     TIMESTAMP                                             NOT NULL,
    closed_at     TIMESTAMP                                             NOT NULL,
    created_at    TIMESTAMP                                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP                                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (class_id, access_code)
);

CREATE TABLE exam_questions
(
    id                   BIGSERIAL PRIMARY KEY,
    exam_id              BIGINT REFERENCES exams (id) ON DELETE CASCADE                 NOT NULL,
    original_question_id BIGINT REFERENCES collection_questions (id) ON DELETE RESTRICT NOT NULL,
    order_num            INTEGER                                                        NOT NULL,
    title                VARCHAR(255)                                                   NOT NULL,
    description          TEXT                                                           NULL,
    check_query          TEXT                                                           NOT NULL,
    check_prompt         TEXT                                                           NOT NULL,
    created_at           TIMESTAMP                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (exam_id, order_num)
);

CREATE TABLE exam_attempts
(
    id              BIGSERIAL PRIMARY KEY,
    exam_id         BIGINT REFERENCES exams (id) ON DELETE CASCADE         NOT NULL,
    class_joinee_id BIGINT REFERENCES class_joinees (id) ON DELETE CASCADE NOT NULL,
    opened_at       TIMESTAMP                                              NULL,
    started_at      TIMESTAMP                                              NULL,
    finished_at     TIMESTAMP                                              NULL,
    database_name   VARCHAR(255)                                           NULL,
    created_at      TIMESTAMP                                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP                                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (exam_id, class_joinee_id)
);

CREATE TABLE exam_submissions
(
    id                  BIGSERIAL PRIMARY KEY,
    exam_question_id    BIGINT REFERENCES exam_questions (id) ON DELETE CASCADE NOT NULL,
    exam_attempt_id     BIGINT REFERENCES exam_attempts (id) ON DELETE CASCADE  NOT NULL,
    answer              TEXT                                                    NOT NULL,
    check_query_passed  BOOLEAN                                                 NULL,
    check_query_at      TIMESTAMP                                               NULL,
    check_prompt_passed BOOLEAN                                                 NULL,
    check_prompt_at     TIMESTAMP                                               NULL,
    created_at          TIMESTAMP                                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP                                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (exam_question_id, exam_attempt_id)
);

-- * auto-update function for updated_at timestamps
CREATE OR REPLACE FUNCTION auto_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- * triggers to automatically update updated_at
CREATE TRIGGER auto_updated_at_users
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_collections
    BEFORE UPDATE
    ON collections
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_collection_questions
    BEFORE UPDATE
    ON collection_questions
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_semesters
    BEFORE UPDATE
    ON semesters
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_classes
    BEFORE UPDATE
    ON classes
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_class_joinees
    BEFORE UPDATE
    ON class_joinees
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_exams
    BEFORE UPDATE
    ON exams
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_exam_questions
    BEFORE UPDATE
    ON exam_questions
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_exam_attempts
    BEFORE UPDATE
    ON exam_attempts
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

CREATE TRIGGER auto_updated_at_exam_submissions
    BEFORE UPDATE
    ON exam_submissions
    FOR EACH ROW
EXECUTE FUNCTION auto_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exam_submissions;
DROP TABLE exam_attempts;
DROP TABLE exam_questions;
DROP TABLE exams;
DROP TABLE class_joinees;
DROP TABLE classes;
DROP TABLE semesters;
DROP TABLE collection_questions;
DROP TABLE collections;
DROP TABLE users;
DROP FUNCTION auto_updated_at;
-- +goose StatementEnd