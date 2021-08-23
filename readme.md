# RestApi PostgreSql, Golang

Такое тестовой, нужно сделать микросервис ленту новостей.
Тип новости: 
• ID
• Title
• Text 
• Date
• Likes
• UserID
Функционал: Каждый пользователь должен просмотреть свою ленту в зависимости от подписок на других пользователей и можно посмотреть кто поставил лайки. Навигация по ленте должна быть курсорной.


# структура бд которую я использовал

CREATE TABLE news (
    id bigserial not null primary key,
    title varchar not null,
    text_news varchar not null,
    date_create timestamp,
    likes JSONB,
    user_id varchar
);

CREATE TABLE users (
    id bigserial not null primary key,
    first_name varchar not null,
    last_name varchar not null,
    password text,
    subscriptions JSONB,
    date_create timestamp,
);

# endpoints:

    /news?id=1&count=5 - получение стартовой страницы новостей

    обязательные параметры:
                id - id пользователя(если добавить авторизацию можно брать из токена)
                count - кол-во записей в ответе

    /fetch?id=1&count=5 - получение последующих страниц при прокрутке

    обязательные параметры:
                id - id пользователя(если добавить авторизацию можно брать из токена)
                count - кол-во записей в ответе

    /likes?id=1 - получение списка пользователей лайкнувших запись

    обязательные параметры:
                id - id новости
                
