-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert_users() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update_users() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete_users() RETURNS TRIGGER AS $$
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
FOR EACH ROW EXECUTE FUNCTION notify_insert_users();

-- Trigger for UPDATE
CREATE TRIGGER users_notify_update
AFTER UPDATE ON public.users
FOR EACH ROW EXECUTE FUNCTION notify_update_users();

-- Trigger for DELETE
CREATE TRIGGER users_notify_delete
AFTER DELETE ON public.users
FOR EACH ROW EXECUTE FUNCTION notify_delete_users();


-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'DELETE', 'table', TG_TABLE_NAME, 'data', row_to_json(OLD))::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS hashtags_notify_insert ON public.hashtags;
DROP TRIGGER IF EXISTS hashtags_notify_update ON public.hashtags;
DROP TRIGGER IF EXISTS hashtags_notify_delete ON public.hashtags;

-- Trigger for INSERT
CREATE TRIGGER hashtags_notify_insert
AFTER INSERT ON public.hashtags
FOR EACH ROW EXECUTE FUNCTION notify_insert_hashtags();

-- Trigger for UPDATE
CREATE TRIGGER hashtags_notify_update
AFTER UPDATE ON public.hashtags
FOR EACH ROW EXECUTE FUNCTION notify_update_hashtags();

-- Trigger for DELETE
CREATE TRIGGER hashtags_notify_delete
AFTER DELETE ON public.hashtags
FOR EACH ROW EXECUTE FUNCTION notify_delete_hashtags();


-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'DELETE', 'table', TG_TABLE_NAME, 'data', row_to_json(OLD))::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS projects_notify_insert ON public.projects;
DROP TRIGGER IF EXISTS projects_notify_update ON public.projects;
DROP TRIGGER IF EXISTS projects_notify_delete ON public.projects;

-- Trigger for INSERT
CREATE TRIGGER projects_notify_insert
AFTER INSERT ON public.projects
FOR EACH ROW EXECUTE FUNCTION notify_insert_projects();

-- Trigger for UPDATE
CREATE TRIGGER projects_notify_update
AFTER UPDATE ON public.projects
FOR EACH ROW EXECUTE FUNCTION notify_update_projects();

-- Trigger for DELETE
CREATE TRIGGER projects_notify_delete
AFTER DELETE ON public.projects
FOR EACH ROW EXECUTE FUNCTION notify_delete_projects();


-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert_user_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update_user_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete_user_projects() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'DELETE', 'table', TG_TABLE_NAME, 'data', row_to_json(OLD))::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS user_projects_notify_insert ON public.user_projects;
DROP TRIGGER IF EXISTS user_projects_notify_update ON public.user_projects;
DROP TRIGGER IF EXISTS user_projects_notify_delete ON public.user_projects;

-- Trigger for INSERT
CREATE TRIGGER user_projects_notify_insert
AFTER INSERT ON public.user_projects
FOR EACH ROW EXECUTE FUNCTION notify_insert_user_projects();

-- Trigger for UPDATE
CREATE TRIGGER user_projects_notify_update
AFTER UPDATE ON public.user_projects
FOR EACH ROW EXECUTE FUNCTION notify_update_user_projects();

-- Trigger for DELETE
CREATE TRIGGER user_projects_notify_delete
AFTER DELETE ON public.user_projects
FOR EACH ROW EXECUTE FUNCTION notify_delete_user_projects();


-- Trigger function for INSERT operation
CREATE OR REPLACE FUNCTION notify_insert_project_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'INSERT', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE operation
CREATE OR REPLACE FUNCTION notify_update_project_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'UPDATE', 'table', TG_TABLE_NAME, 'data', row_to_json(NEW))::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for DELETE operation
CREATE OR REPLACE FUNCTION notify_delete_project_hashtags() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('crud_operations', json_build_object('operation', 'DELETE', 'table', TG_TABLE_NAME, 'data', row_to_json(OLD))::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS project_hashtags_notify_insert ON public.project_hashtags;
DROP TRIGGER IF EXISTS project_hashtags_notify_update ON public.project_hashtags;
DROP TRIGGER IF EXISTS project_hashtags_notify_delete ON public.project_hashtags;

-- Trigger for INSERT
CREATE TRIGGER project_hashtags_notify_insert
AFTER INSERT ON public.project_hashtags
FOR EACH ROW EXECUTE FUNCTION notify_insert_project_hashtags();

-- Trigger for UPDATE
CREATE TRIGGER project_hashtags_notify_update
AFTER UPDATE ON public.project_hashtags
FOR EACH ROW EXECUTE FUNCTION notify_update_project_hashtags();

-- Trigger for DELETE
CREATE TRIGGER project_hashtags_notify_delete
AFTER DELETE ON public.project_hashtags
FOR EACH ROW EXECUTE FUNCTION notify_delete_project_hashtags();

