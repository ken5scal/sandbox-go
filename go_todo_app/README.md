

# ?
 * errgroup.Group
 * type
 * store/repository.goの各TX Interfaceはjmoiron/sqlx前提で使ってるからそうなってるのか

# memo
 * docker compose build --no-cache
 * docker compose up
 * openssl genrsa 4096 > ./auth/cert/secret.pem
 * openssl rsa -pubout < ./auth/cert/secret.pem > ./auth/cert/public.pem
 * curl -i -XPOST localhost:18000/register -d @./handler/testdata/register_user/ok_req.json.golden
 * curl -i -XPOST localhost:18000/register -d '{"name": "admin_user", "password":"dQ7fLPtzFl7UhqKc68lZO", "role":"admin"}'
 * curl XPOST localhost:18000/login -d '{"user_name": "john", "password":"test"}' | jq -r .access_token | pbcopy && export TODO_TOKEN=$(pbpaste)
 * curl -i XPOST localhost:18000/login -d '{"user_name": "admin_user", "password":"dQ7fLPtzFl7UhqKc68lZO"}'
 * curl -i -XPOST -H "Authorization: Bearer $TODO_TOKEN" localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden -v
 * curl -i -XGET -H "Authorization: Bearer $TODO_TOKEN" localhost:18000/tasks


# 元ネタ
 * [公式コード](https://github.com/budougumi0617/go_todo_app)