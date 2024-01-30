## Consumer Benchmark

Данный тест можно выполнять на любом из брокеров

Бенчмарк консюмера запускается скриптом `kafka-consumer-perf-test.sh`. Рассмотрим его основные конфиги:

```bash
kafka-consumer-perf-test --help
```

- **bootstrap-server** указывает адрес брокера к которому будет подключаться консюмер
- **messages** обозначает количество сообщений которые должны быть зафетчены консюмером
- **topic** указывает название топика из которого бенчмарк будет читать сообщения
- **print-metrics** выводит детальные метрики в конце работы бенчмарка

Создадим топик для того, чтобы провести тест консюмера:

```bash
kafka-topics --bootstrap-server localhost:9092 --topic consumer-benchmark --replication-factor 3 --partitions 1 --create
```

Запишем в него набор сообщений, воспользовавшись бенчмарком продюссера:

```bash
kafka-producer-perf-test --producer-props bootstrap.servers=localhost:9092 --topic consumer-benchmark --throughput -1 --num-records 1000000 --record-size 100
```

Теперь запустим бенчмарк консюмера:

```bash
kafka-consumer-perf-test --bootstrap-server localhost:9092 --topic consumer-benchmark --messages 1000000 --print-metrics
```

Получаются следующие метрики

```
start.time
2024-01-30 17:45:07:644

end.time
2024-01-30 17:45:12:950

data.consumed.in.MB
95.3674

MB.sec
17.9735

data.consumed.in.nMsg
1000000

nMsg.sec
188465.8877

rebalance.time.ms
4463

fetch.time.ms
843

fetch.MB.sec
113.1286

fetch.nMsg.sec
1186239.6204
```

Как только тест завершится вы увидите набор метрик в терминале, основные из которых:
- **start.time, end.time** — время начала и конца теста
- **data.consumed.in.MB** — общий размер всех считанных сообщений
- **MB.sec** — пропускная способность консюмера в МБ/сек (объем считанных данных за секунду)
- **data.consumed.in.nMsg** — общее количество считанных сообщений
- **nMsg.sec** — пропускная способность консюмера в сообщениях/сек

В приведенном примере мы видим, что консюмер вычитал 1,000,000 сообщений весом ±95 МБ данных с пропускной способностью ±17 МБ/сек за 5 сек.

Попробуем создать топик с несколькими партициями и сравнить результат:

```bash
kafka-topics --bootstrap-server localhost:9092 --topic consumer-benchmark-multipartition --replication-factor 3 --partitions 3 --create
```

```bash
kafka-producer-perf-test --producer-props bootstrap.servers=localhost:9092 --topic consumer-benchmark-multipartition --throughput -1 --num-records 1000000 --record-size 100
```

```bash
kafka-consumer-perf-test --bootstrap-server localhost:9092 --topic consumer-benchmark-multipartition --messages 1000000
```

```
start.time
2024-01-30 17:52:32:569

end.time
2024-01-30 17:52:36:647

data.consumed.in.MB
95.3674, 

MB.sec
23.3858

data.consumed.in.nMsg
1000000

nMsg.sec
245218.2442

rebalance.time.ms
3388

fetch.time.ms
690

fetch.MB.sec
138.2137

fetch.nMsg.sec
1449275.3623
```

В приведенном примере мы видим, что консюмер вычитал 1,000,000 сообщений весом ±95 МБ данных с пропускной способностью ±23 МБ/сек за 4 сек.
