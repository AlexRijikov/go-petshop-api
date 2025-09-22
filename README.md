# go-petshop-api
golang petshop api
// Архітектура проєкту 
       (Version-1.0)
/cmd
    /server         → запуск API
/internal
    /config         → конфіг (env, підключення до БД)
    /models         → опис моделей (User, Product, Order)
    /repositories   → робота з БД
    /services       → бізнес-логіка
    /handlers       → REST-ендпоінти
    /middleware     → JWT, логування
/pkg
    /utils          → допоміжні функції
