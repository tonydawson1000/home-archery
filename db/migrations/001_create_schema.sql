CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE archers (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name  text NOT NULL UNIQUE,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    archer_id       uuid NOT NULL REFERENCES archers (id),
    started_at      timestamptz NOT NULL DEFAULT now(),
    completed_at    timestamptz,
    status          text NOT NULL DEFAULT 'in_progress'
                    CHECK (status IN ('in_progress', 'completed')),
    arrows_per_end  smallint NOT NULL DEFAULT 6
                    CHECK (arrows_per_end = 6),
    arrow_count     integer NOT NULL DEFAULT 0 CHECK (arrow_count >= 0),
    total_score     integer NOT NULL DEFAULT 0 CHECK (total_score >= 0),
    hit_count       integer NOT NULL DEFAULT 0,
    gold_count      integer NOT NULL DEFAULT 0,
    x_count         integer NOT NULL DEFAULT 0,
    CONSTRAINT sessions_complete_consistency CHECK (
        (status = 'in_progress' AND completed_at IS NULL)
        OR (status = 'completed' AND completed_at IS NOT NULL AND arrow_count > 0)
    )
);

CREATE TABLE ends (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   uuid NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    end_number   integer NOT NULL CHECK (end_number >= 1),
    UNIQUE (session_id, end_number)
);

CREATE TABLE arrows (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    end_id        uuid NOT NULL REFERENCES ends (id) ON DELETE CASCADE,
    arrow_number  smallint NOT NULL CHECK (arrow_number BETWEEN 1 AND 6),
    score_code    text NOT NULL
                  CHECK (score_code IN ('X','10','9','8','7','6','5','4','3','2','1','M')),
    score_value   smallint NOT NULL CHECK (score_value BETWEEN 0 AND 10),
    recorded_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (end_id, arrow_number)
);

CREATE INDEX idx_sessions_archer_started ON sessions (archer_id, started_at DESC);
CREATE INDEX idx_sessions_archer_pb
    ON sessions (archer_id, arrow_count, total_score DESC)
    WHERE status = 'completed';
