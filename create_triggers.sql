-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'DELETE', 'table', TG_TABLE_NAME, 'data', row_to_json(OLD))::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS users_notify_insert ON public.users;
DROP TRIGGER IF EXISTS users_notify_update ON public.users;
DROP TRIGGER IF EXISTS users_notify_delete ON public.users;

-- Trigger for INSERT
CREATE TRIGGER users_notify_insert
AFTER INSERT ON public.users
FOR EACH ROW EXECUTE FUNCTION notify_insert();

-- Trigger for UPDATE
CREATE TRIGGER users_notify_update
AFTER UPDATE ON public.users
FOR EACH ROW EXECUTE FUNCTION notify_update();

-- Trigger for DELETE
CREATE TRIGGER users_notify_delete
AFTER DELETE ON public.users
FOR EACH ROW EXECUTE FUNCTION notify_delete();
