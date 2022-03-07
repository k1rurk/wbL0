package main

import (
	"wb_l0/cache"
	"wb_l0/database"
	"wb_l0/httpServer"
	"wb_l0/natJetStream/sub"
)

func main() {
	//Открываем базу данных
	db := database.Open()

	//Вытаскиваем все заказы из базы данных
	orders := database.LoadDataFromDb(db)

	//Запихиваем их в кэш
	c := cache.NewCache(orders)

	//Подписываемся jet stream для получения данных
	natJetStream.Sub(c, db)
	//defer subscript.Unsubscribe()

	//Запускаем сервер
	server := httpServer.InitServer(c)
	server.StartServer()
}
