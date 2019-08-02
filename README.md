# ispjournalctl
## Описание
Утилита для чтения журналов `isp` модулей.
Позволяет преобразовать внутренний бинарный формат записей в человекочитаемые форматы для дальнейшей обработки.(`csv`, `json`)
## Требования
* Linux
## Установка
```sh
yum install ispjournalctl
```
## Конфигурация
Внешняя конфигурация отсутствует
## Использование
### Чтение локальных файлов журналов
```sh
ispjournalctl read -h
```
```sh
Read isp journal file

Usage:
  ispjournalctl read     [flags]

Flags:
      --event strings   filtered events, format: [--event='event1' --event='event2'], empty: show all
      --file string     source file to read
      --gz              source file is gzipped
  -h, --help            help for read
      --level strings   filtered log levels, format: [--level='OK' --level='WARN', --level='ERROR'], empty: show all
  -n, --n int           log entries count from start, defautl: read all
  -o, --out string      output format in csv with ';' or json, example: --out='csv' (default "csv")
      --since string    since time in format 2018-06-15 [08:15:00]
      --until string    until time in format 2018-06-15 [08:15:00]
```
Преобразует данные полученные из `stdin` или файла указанного через флаг `--file`.
Позволяет фильтровать записи логов по событиям, уровням логирования, времени.
Результаты записываются в `stdout` построчно, в указнном формате.
### Чтение файлов журналов из isp-journal-service
```sh
ispjournalctl search -h
```
```sh
Search isp journal file

UUsage:
   ispjournalctl search [flags]
 
 Flags:
   -g, --gate string     gate to isp-journal-service in format '127.0.0.0:0000'
       --event strings   filtered events, format: [--event='event1' --event='event2'], empty: show all
   -h, --help            help for search
       --host strings    filtered host, format: [--host='host1' --host='host2'], empty: show all
       --level strings   filtered log levels, format: [--level='OK' --level='WARN', --level='ERROR'], empty: show all
       --module string   module name
   -n, --n int           log entries count from start, default: read all (default -1)
   -o, --out string      output format in csv with ';' or json, example: --out='csv' (default "csv")
       --since string    since time in format 2018-06-15 [08:15:00]
       --until string    until time in format 2018-06-15 [08:15:00]

```
Читает логи указанного модуля `--module` подключаясь с помощью `--gate` к `isp-journal-service`
Позволяет фильтровать записи логов по событиям, уровням логирования, времени.
Результаты записываются в `stdout` построчно, в указнном формате.

#### Пример
```sh
ispjournalctl read --gz  --out json --level OK --level ERROR --file "10.15.27.48__2019-07-10T06-59-45.655.log.gz"
```
```json
{"moduleName":"example","host":"111.15.29.48","event":"test","level":"OK","time":"2019-07-09T11:08:32.542+00:00","request":"{\n\t\"objId\": \"e5d7c9ae-93ef-4a12-b35c-e4199d6fd3ec\"\n}","response":"{\"timestamp\":1562670512,\"random\":521507349,\"secureHash\":\"3d439713b59e103e0ccd31a7cb9de12f3433fd814f2ae75632261e485b42b4c0\",\"code\":[14],\"desc\":\"OK\"}"}
{"moduleName":"example","host":"111.15.29.48","event":"test","level":"OK","time":"2019-07-09T11:08:47.936+00:00","request":"{\n\t\"objId\": \"92cb4330-ad2c-11e9-a2a3-2a2ae2dbcce4\"\n}","response":"{\"timestamp\":1562670527,\"random\":719003178,\"secureHash\":\"11094cfa98f983b12150c2936480160f692b21986afabcea06a70074a562e0c8\",\"code\":[14],\"desc\":\"OK\"}"}
```
```sh
ispjournalctl read --gz --file "10.15.27.48__2019-07-10T06-59-45.655.log.gz"
```
```csv
module_name;host;event;level;time;request;response;error_text
example;111.15.29.48;test;OK;2019-07-09T11:08:32.542+00:00;"{
        ""objId"": ""e5d7c9ae-93ef-4a12-b35c-e4199d6fd3ec""
}";"{""timestamp"":1562670512,""random"":521507349,""secureHash"":""3d439713b59e103e0ccd31a7cb9de12f3433fd814f2ae75632261e485b42b4c0"",""code"":[14],""desc"":""OK""}";
example;111.15.29.48;test;OK;2019-07-09T11:08:47.936+00:00;"{
        ""objId"": ""92cb4330-ad2c-11e9-a2a3-2a2ae2dbcce4""
}";"{""timestamp"":1562670527,""random"":719003178,""secureHash"":""11094cfa98f983b12150c2936480160f692b21986afabcea06a70074a562e0c8"",""code"":[14],""desc"":""OK""";
```


```sh
ispjournalctl search --gate '127.0.0.1:0000' --module 'example' --out 'csv'
```
```csv
example;111.15.29.48;test;OK;2019-07-09T11:08:32.542+00:00;"{
        ""objId"": ""e5d7c9ae-93ef-4a12-b35c-e4199d6fd3ec""
}";"{""timestamp"":1562670512,""random"":521507349,""secureHash"":""3d439713b59e103e0ccd31a7cb9de12f3433fd814f2ae75632261e485b42b4c0"",""code"":[14],""desc"":""OK""}";
example;111.15.29.48;test;OK;2019-07-09T11:08:47.936+00:00;"{
        ""objId"": ""92cb4330-ad2c-11e9-a2a3-2a2ae2dbcce4""
}";"{""timestamp"":1562670527,""random"":719003178,""secureHash"":""11094cfa98f983b12150c2936480160f692b21986afabcea06a70074a562e0c8"",""code"":[14],""desc"":""OK""";
```