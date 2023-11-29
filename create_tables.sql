-- Create Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    created_at TIMESTAMP
);

-- Create Hashtags table
CREATE TABLE hashtags (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    created_at TIMESTAMP
);

-- Create Projects table
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    slug VARCHAR,
    description TEXT,
    created_at TIMESTAMP
);

-- Create Project_Hashtags table with foreign key references
CREATE TABLE project_hashtags (
    hashtag_id INT REFERENCES hashtags(id),
    project_id INT REFERENCES projects(id),
    PRIMARY KEY (hashtag_id, project_id)
);

-- Create User_Projects table with foreign key references
CREATE TABLE user_projects (
    project_id INT REFERENCES projects(id),
    user_id INT REFERENCES users(id),
    PRIMARY KEY (project_id, user_id)
);
