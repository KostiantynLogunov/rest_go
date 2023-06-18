# rest_api_tutorial

# user-service

# REST API

GET /users -- list of users -- 200, 404, 500
GET /users/:id -- user by id -- 200, 404, 500
POST /users-- create a new user -- 201, 4xx, Header Location: url
PUT /users/:id -- fully update user by id -- 200, 204, 404, 500
PATCH /users/:id -- partially update user by id -- 200, 204, 404, 500
DELETE /users/:id -- remove user by id --  204, 404, 500

204 - no content
200-ok