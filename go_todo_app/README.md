

# ?
 * errgroup.Group
 * type
 * store/repository.goの各TX Interfaceはjmoiron/sqlx前提で使ってるからそうなってるのか

# memo
 * docker compose build --no-cache
 * docker compose up
 * curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden -v
 * curl -i -XGET localhost:18000/tasks
 * curl -i -XPOST localhost:18000/register -d @./handler/testdata/register_user/ok_req.json.golden -v
 * curl -i XPOST localhost:18000/login -d '{"user_name": "john", "password":"dQ7fLPtzFl7UhqKc68lZO"}'
 * openssl genrsa 4096 > secret.pem
 * openssl rsa -pubout < secret.pem > public.pem


# 元ネタ
 * [公式コード](https://github.com/budougumi0617/go_todo_app)