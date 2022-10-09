# Currency exchange 

1) GET http://10.10.10.1:10000/rates?onDate=07.10.2022 - получение курсов валют, на дату из цб
2) GET http://10.10.10.1:10000/exchange?from=USD&to=RUB&onDate=07.10.2022 получение курса на дату из базы

## Работа с курсами
ЦБ присылает 34 пары курсов валют. Из них получаются обратные курсы путем 1 / exchangeRate<br>
После происходит триангуляция курсов путем sql запроса вида:

```sql
insert into currency_rates (currency_from, currency_to, created_at, on_date, exchange_rate)
select pair.currency_from, pair.currency_to, pair.created_at, pair.on_date, pair.exchange_rate
from (
         select f.currency_from, t.currency_from as currency_to, f.created_at, f.on_date, (f.exchange_rate / t.exchange_rate) as exchange_rate
         from currency_rates f, currency_rates t
         where f.on_date = t.on_date
           and f.currency_to = t.currency_to
     ) pair
         LEFT OUTER JOIN currency_rates cr
                         ON (
                                     pair.on_date = cr.on_date
                                 AND pair.currency_from = cr.currency_from
                                 AND pair.currency_to = cr.currency_to
                             )
where pair.on_date = '2022-10-07'
group by pair.currency_from, pair.currency_to
```
[Источник](http://www.dpxo.net/articles/fx_rate_triangulation_sql.html)

## todo:
1) обмен валют
2) обработка ошибок
3) логгирование
4) валидация реквеста
5) докер
6) постгрес
