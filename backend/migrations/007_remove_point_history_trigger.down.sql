-- Migration 007 Down: Restore automatic point history trigger
-- (Copied from 002_create_triggers_and_functions.up.sql)

CREATE OR REPLACE FUNCTION record_point_history()
RETURNS TRIGGER AS $$
DECLARE
    point_type VARCHAR(50);
    point_amount BIGINT;
BEGIN
    -- 変更内容に基づいてポイント種別と金額を判定
    IF NEW.total_points > OLD.total_points THEN
        -- ポイント付与
        point_type := 'grant';
        point_amount := NEW.total_points - OLD.total_points;
    ELSIF NEW.reserved_points > OLD.reserved_points THEN
        -- 入札時の予約
        point_type := 'reserve';
        point_amount := -(NEW.reserved_points - OLD.reserved_points);
    ELSIF NEW.reserved_points < OLD.reserved_points AND NEW.available_points > OLD.available_points THEN
        -- 予約解放
        point_type := 'release';
        point_amount := NEW.available_points - OLD.available_points;
    ELSIF NEW.reserved_points < OLD.reserved_points AND NEW.available_points = OLD.available_points THEN
        -- 落札時の消費
        point_type := 'consume';
        point_amount := -(OLD.reserved_points - NEW.reserved_points);
    ELSE
        -- その他の変更（通常は発生しない）
        RETURN NEW;
    END IF;

    -- 履歴レコードの挿入
    INSERT INTO point_history (
        bidder_id,
        amount,
        type,
        balance_before,
        balance_after,
        reserved_before,
        reserved_after,
        total_before,
        total_after
    ) VALUES (
        NEW.bidder_id,
        point_amount,
        point_type,
        OLD.available_points,
        NEW.available_points,
        OLD.reserved_points,
        NEW.reserved_points,
        OLD.total_points,
        NEW.total_points
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_record_point_history AFTER UPDATE ON bidder_points
    FOR EACH ROW
    WHEN (OLD.* IS DISTINCT FROM NEW.*)
    EXECUTE FUNCTION record_point_history();
