# Pinger

## Функционал

Фронтенд отрисовывает табличку, в которую помещает данные полученные из бэкэнда.
Бэк же берёт данные из базы данных(postgreSQL)

Сервис pinger после запуска начинает смотреть запущенные докер контейнеры и пингует их(проверяет работает ли).

После пинга каждого контейнера, сервис отправляет в БД: ip, длительность пинга в мс, время последнего успешного пинга.

## Запуск

Для начала нунжо склонировать репозиторий:

```githubexpressionlanguage
git clone https://github.com/k1v4/Pinger
```

Предварительно стоит позаботиться, чтобы порты 8080, 5432 и 3000 были свободны.  
Далее стоит перейти в папку репозитория и использовать команду для старта docker compose(потребуется Docker для запуска системы):

```shell
  docker-compose up
```
 
Если же у вас установлена утилита make, то также находясь в папке проекта, можно просто ввести команду:

```shell
  make
```

После выполнения команды по запуску стоит немного подождать, чтобы все сервисы запустились

Для получения результатов стоит перейти на [http://localhost:3000 ](http://localhost:3000)