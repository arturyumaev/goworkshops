## Producer Benchmark

Запускаем кластер из 3 брокеров и одного ZooKeeper командой `make`

Заходим в консоль одного из брокеров, все скрипты для управления брокером уже предустановлены компанией confluent в `/bin` и доступы в любой директории без префикса `*.sh`

```bash
ls -Hla /bin | grep "kafka"
```

Скрипт бенчмарка продюсера называется kafka-producer-perf-test, давайте ознакомимся с набором его опций

```bash
kafka-producer-perf-test

usage: producer-performance [-h] --topic TOPIC --num-records NUM-RECORDS [--payload-delimiter PAYLOAD-DELIMITER] --throughput THROUGHPUT [--producer-props PROP-NAME=PROP-VALUE [PROP-NAME=PROP-VALUE ...]]
                            [--producer.config CONFIG-FILE] [--print-metrics] [--transactional-id TRANSACTIONAL-ID] [--transaction-duration-ms TRANSACTION-DURATION] (--record-size RECORD-SIZE |
                            --payload-file PAYLOAD-FILE)
```


Нас интересуют следующие конфиги:
- **producer-props** позволяет нам передавать любые опции инстансу продюсера
- **topic** указывает название топика в который бенчмарк будет писать данные
- **num-records** устанавливает лимит на общее число отправленных сообщений
- **throughput** устанавливает лимит на число отправляемых сообщений в секунду
- **record-size** конфигурирует размер одного сообщения в байтах
- **print-metrics** указывает бенчмарку вывести на экран все финальные метрики продюсера

---

В первую очередь нам нужно создать топик для теста. Создадим топик с одной партицией и RF 3

```bash
kafka-topics --bootstrap-server localhost:9092 --topic producer-benchmark --replication-factor 3 --partitions 1 --create
```

---

Запустим бенчмарк продюсера и укажем ему записать 50000 сообщений, с пропускной способностью 1000 сообщений в секунду  размером 100 байт:

```bash
kafka-producer-perf-test --producer-props bootstrap.servers=localhost:9092 --topic producer-benchmark --throughput 1000 --num-records 50000 --record-size 100 --print-metrics
```

```
...

50000 records sent, 1000.020000 records/sec (0.10 MB/sec), 2.09 ms avg latency, 332.00 ms max latency, 1 ms 50th, 3 ms 95th, 28 ms 99th, 94 ms 99.9th.

...
```

Во время своей работы бенчмарк выводит на экран метрики производительности продюсера — число посланных сообщений, пропускную способность, среднюю и максимальную задержку.

Из приведенного выше результата мы можем сделать вывод, что мы можем записать в 1 партицию 100 KB данных в секунду со средней задержкой в ±2ms, задержка 99 персентиля равна 28ms.

---

Посмотрим, какую максимальную пропускную способность мы можем ожидать от 1 партиции в нашем кластере при размере одного сообщения в 100 байт. Для этого передадим значение -1 в параметр --throughput, тем самым убрав ограничение на количество сообщений в секунду, а так же выставим --num-records в 500000.

```bash
kafka-producer-perf-test --producer-props bootstrap.servers=localhost:9092 --topic producer-benchmark --throughput -1 --num-records 500000 --record-size 100
```

```
500000 records sent, 132345.156167 records/sec (12.62 MB/sec), 227.47 ms avg latency, 871.00 ms max latency, 214 ms 50th, 801 ms 95th, 863 ms 99th, 869 ms 99.9th
```

По результатам можно установить, что продюсер с дефолтными настройками может записать ±12 МБ/сек данных со средней задержкой в 227 мс.

---

По умолчанию продюсер Кафки использует acks=1, т.е. ожидает ответа успешной записи от лидера партиции. Выставим это значение в 0, таким образом сказав продюсеру не ждать никаких ответов

```
kafka-producer-perf-test --producer-props bootstrap.servers=localhost:9092 acks=0 --topic producer-benchmark --throughput -1 --num-records 500000 --record-size 100
```

```
500000 records sent, 266951.414842 records/sec (25.46 MB/sec), 98.83 ms avg latency, 332.00 ms max latency, 89 ms 50th, 191 ms 95th, 303 ms 99th, 331 ms 99.9th
```

Видим, что производительность увеличилась и мы можем писать в среднем 25 MB/sec со средней задержкой 98ms

---

Попробуй выставить `acks=-1` и посмотри как это повлияет на производительность
