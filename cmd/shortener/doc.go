// Сервис сокращения ссылок go_link_shortener
//
// При обращении к сервису с передачей URL, создает уникальный идентификатор
// на основе алгоритма base62 и сохраняет в одном из трех вариантов хранилища
//
//  1. Slice в памяти;
//  2. Запись в файл;
//  3. Postgresql DB;
//
// Логика взаимодействия с хранилищем располгагается в одном из 3х соответствующих репозиториев
// после обращения из сервисного слоя.
//
// # Сборка приложения происходит соответствующей командой:
//
//	go build
//
// При сборке допустима передача флагов компиляции для сохранения в переменных
// и последующего отображения в stdout при запуске приложения.
//
// Version information variables. Initialized durign the build process.
// For example, use next command for build app shortener
//
//	buildVersion
//	buildDate
//	buildCommit
//
// # Пример вызова команды сборки с передачей соответствующих значений переменных в флагах компиляции:
//
//	go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildCommit=$(git rev-parse HEAD~1)'"
//
// # Команда запуска приложения из текущей директории:
//
//	./shortener
//
// Поддерживается передача ряда конфигурационных флагов для запуска приложения с опеределенными параметрами:
//
//	-a // значениe хоста для вызова приложения.
//	-a=localhost:8081
//
//	-b // значениe для подстановки хоста к сокращенному URL
//	-b=localhost:9090
//
//	-l // уровень логирования
//	-l=info
//
//	-f // полное название файла хранилища сокращенных URL, содержащее путь к нему
//	-f=/tmp/file_storage_name.json
//
//	-d // строка DSN конфигурации базы данных
//	-d=host=localhost port=5432 user=postgres password=root dbname=urls sslmode=disable
//
// # Пример кманды запуска прилжения с передачей конфигурационных параметров при запуске:
//
//	./shortener -a=localhost:8081 -b=localhost:9090 -l=info -d="host=localhost port=5432 user=postgres password=root dbname=urls sslmode=disable"
package main
