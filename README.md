# Go-balancer

Балансировщик нагрузки, работающий с базой данных и собирающий докер-контейнер.

Запуск:
`make up`

Реализовано:

1. Получение данных конфигурации с .json и .env файлов. 
   * config.json - порт сервера и адреса бэкендов
   * bucketConfig - конфигурация количества и темпа накопления бакетов для разных ip адресов.
   * config.env - конфигурация базы данных
2. Распределение запросов по бекендам методом round-robin
3. Все параллельно используемые данные защищены мьютексом
4. Ошибки логируются стандартным логгером.
5. Все подключения логируются через middleware и стандартный логер в нем
6. При недоступности сервера - перебрасывает на ближайший (по id) доступный
7. Проверка жизни серверов раз в 5 секунд
8. Каждому подключившемуся клиенту выдается свой бакет токенов, если его нет в базе данных или в конфигурации, то присваиваются стандартные значения
9. Все конфигурации для клиентов сохраняются в базу данных с тремя полями: ip, capacity, rate
10. При каждом новом получении страницы от пользователя у него обновляются его значение токенов
11. Сборка в докер контейнер с тремя сервисами:
    * database - база данных
    * app - балансировщик
    * test_utils - десять серверов на локалхосте для проверки
12. Есть коллекция postman для нагрузки и проверки параллельных запросов
13. Makefile для запуска одной командой