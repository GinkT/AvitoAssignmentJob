### Деплой проекта

Для запуска проекта созданы docker-compose и docker файлы.  
Достаточно прописать команду и подтянуть образы:    

`docker-compose up`

Сервер слушает запросы по порту **8181**. Также, в compose подключен    
образ adminer, который работает на порту **8080**. С помощью него можно     
удобно проверять наполнение БД. Данные для входа в админер можно взять из   
main.go файла(константы).

БД инициализируется с файла _init.sql_. Создаётся 2 таблицы: _users, transactions_

### Спецификация HTTP API:
 
Реализованы все требуемые в тех. задании методы(и дополнительные):

**balance**
`http://localhost:8181/balance?id=1`    `http://localhost:8181/balance?id=1&currency=USD`

**payment**
`http://localhost:8181/payment?id=1&amount=130`

**withdraw**
`http://localhost:8181/withdraw?id=1&amount=130`

**transfer**
`http://localhost:8181/transfer?fromId=1&toId=4&amount=10`

**history**
`http://localhost:8181/history?id=1` `http://localhost:8181/history?id=1&sortBy=amount&orderBy=desc`

В качестве времени хранились значения в формате Unix timestamp. Сделано это было для упрощения и наглядности.   
Поля _sortBy_ и _orderBy_ опциональны. Как и поле _currency_ в методе balance.

### Тесты:

Функции выполняющие запросы к БД снабжены Unit тестами.


_Дата начала выполнения задания: 18.09.2020_ 
_Дата отправки заявки с выполненным заданием: 20.09.2020_
