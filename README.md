# Order service

Данный сервис создан для управления некоторыми абстрактными заказами. Он выпоняет CRUD фукнционал для заказов. Для общения с другими возможными сервисами используется gRPC, но также создан и gateway для использования http запросов.

## Функционал

- Общение с сервисом с использованием grpc, а также добавлен gateway для общения с сервисом посредствам http
- Используется postgreSQL, который поднимается в контейнере
- Сам сервис упакован в docker
- Использован nginx для балансировки приходящих запросов: запросы уходят обрабатываться на 3 сервера по очереди
- Написан docker-compose файл, который включает в себя три копии сервиса, базу данных postgreSQL, а также балансировщик nginx
- Построен простой pipeline для CI/CD на платформе gitlab
