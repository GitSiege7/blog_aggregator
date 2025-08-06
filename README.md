# Blog Aggre-Gator

### Requirements
Compile Requirements:
- Go: 1.24+
- Postgres: 16.9+

### Installation
1. install go
2. install postgres: `sudo apt install postgresql postgresql-contrib postgresql-server`
3. `git clone https://github.com/GitSiege7/blog_aggregator`
4. `cd ./blog_aggregator/`
5. `go install`
6. Create user `postgres`: set `sudo passwd postgres` to `postgres`
7. Create database: 
    - `sudo service postgresql start`
    - `sudo -u postgres psql`
    - `create database gator;`
    - `\c gator`
    - `alter user postgres password 'postgres';`
    - `exit`
8. Migrate database:
    - change dir to `schema`: `cd .../blog_aggregator/sql/schema`
    - migrate up: `goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up`
    - `./run.sh`

### Usage
On first use, you'll want to register a user, and add feeds.

Command Guideline:
- `register {name}`: registers a user into the database
- `login {name}`: logs user into system as current user
- `users`: lists users
- `agg {time interval}`: collects updates on feeds followed by current user, intended to be run in separate terminal in background.
- `addfeed {name} {url}`: adds a feed to database and follows for current user
- `feeds`: lists all feeds and users following
- `follow {url}`: follows a feed for current user
- `following`: lists feeds followed by current user
- `unfollow {url}`: unfollows feed for current user
- `browse {limit}`: displays `limit (def=2)` number of posts from feeds followed by user