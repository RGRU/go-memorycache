# Go-memorycache [![Build Status](https://travis-ci.org/RGRU/go-memorycache.svg?branch=master)](https://travis-ci.org/RGRU/go-memorycache)

Менеджер кеша в памяти на Golang, хранилище данных в формате ключ/значение


## Как установить?

	go get github.com/RGRU/go-memorycache


## Как использовать?

1. Необходимо импортировать пакет

	import (
		memorycache "github.com/RGRU/go-memorycache"
	)

2. Инициализировать кеш
	Создаем кеш с временем жизни по-умолчанию равным 5 минут и удалением просроченного кеша каждые 10 минут
	cache := memorycache.New(5*time.Minute, 10*time.Minute)


3. Использовать

	// Установить кеш с ключем "myKey" и временем жизни 5 минут
	cache.Set("myKey", "My value", 5 * time.Minute)

	// Получить кеш с ключем "myKey"
	cache.Get("myKey")

Полный список методов

| Метод | Описание  | Пример использования |
|:--        |:--   |:--          |
| New | Инициализация кеша. Принимает два аргументы defaultExpiration и cleanupInterval. defaultExpiration - время жизни кеша по-умолчанию, если установлено значение <= 0, время жизни кеша  бессрочно. cleanupInterval - интервал между удалением просроченного кеша. Внимание! не рекомендуется устанавливать небольшие значения для этого аргумента. Если установлено значение <= 0, очистка и удаление просроченного кеша не происходит | cache := memorycache.New(5*time.Minute, 10*time.Minute) |
| Set | Установка кеша. Принимает три аргумента: название ключа (должно быть уникально в рамках одной инициализации), значение и время жизни. Возвращает ошибку в случае если установка кеша не возможна | cache.Set("myKey", "My value", 5 * time.Minute) |
| Get | Получение кеша по ключу. Возвращает значение кеша и true, в случае успешного получения или false если значение не найдено или устарело | value, found := cache.Get("myKey") |
|GetLikeKey | Получение списка кешей с фильтром по имени ключа. При поиске можно использовать специалиный символ %, который работает по аналогии с LIKE в Mysql | values, found := AppCache.GetLikeKey("findkey%") |
|GetCount | Получение кол-ва ключей | count := cache.GetCount() |
|GetItem | Получение экземпляра кеша. Возвращает Item: Value (значение), Duration (время жизни), Created (время создания), Expiration (продолжительность жизни) | value, found := AppCache.GetItem(testKey) |
| Delete | Удаление кеша по ключу. Возвращает ошибку в случае если удаление не возможно | cache.Delete("myKey") |
|Rename | Переименование ключа | error := AppCache.Rename("keyOldName", "keyNewName") |
|Copy | Копирование значения ключа | error := AppCache.Copy("oneKey", "oneKeyCopy") |
| Exist | Проверить кеша на существование. Возвращает true если кеш существует. | exist := cache.Exist("myKey") |
| Expire | Проверка кеша на истечение срока жизни. Возвращает true если время жизни кеша истекло | expire := cache.Exist("myKey") |
| FlushAll | Очистка всех данных. Возвращает ошибку в случае если очистка не возможна | flush := cache.FlushAll() |
