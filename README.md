# mister-aladin-articles

I use existing libs :
 - Echo Router
 - Viper, for config management
 - sqlx, for database connection

I'am build a service to manage data articles using golang, using postgresql for database, and redis for caching. service has database migration.
I created 5 endpoints :
1. [POST] /articles to post a new article
2. [GET] /articles to get the list of articles. Sort the articles by newest first with optional query parameters:
    - a. query: to search keywords in the article title and body
    - b. author: filter by author name.
3. [GET] /articles/:id to get a specific article by :id
4. [PUT] /articles/:id to update an article by :id
5. [DELETE] /articles/:id to delete an article by :id

# Setup after cloning the repo:
- $ cd mister-aladin-articles
- $ go get all
- $ go mod tidy

# connfiguration environment
do this following actions untuk set up your configuration :
- copy config/config-dev.yaml.example ke config/config-dev.yaml
- complete the required credentials according to the existing format

# Database migration :
I use PostgreSQL for Database.
you can create database tables by migration, but before that you have to create a new database on your RDBMS.
- $ migrate -database "mysql://user:password@tcp(host:port)/database" -path pkg/database/migrations up

# Run service :
You can run the service by using the following command, after that the service is ready to use.
- $ go run main.go

# Use the service
you can use services with consumption to API endpoints that are already created, you can use the postman tool. to make it easier to use I've created a workspace for it, and it's ready to use
https://www.postman.com/lively-comet-875863/workspace/mister-aladin-test-workspace/collection/18402968-42a65d81-520a-471e-b619-75a9aa10fe7c?action=share&creator=18402968

# Thank You
if there are problems or want to know more information about me, please contact via linkedin via the following link https://www.linkedin.com/in/adi-khoiron-hasan or by sending an email to adikhoironhasan@gmail.com
