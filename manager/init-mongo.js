rs.initiate({
    _id: "rs0",
    members: [
        { _id: 0, host: "mongo-primary:27017" },
        { _id: 1, host: "mongo-secondary1:27017" },
        { _id: 2, host: "mongo-secondary2:27017" }
    ],
    settings: {
        getLastErrorDefaults: {
            w: 3,
            j: true,
            wtimeout: 5000
        }
    }
})

rs.status()


// Параметр w запрашивает подтверждение того,
// что операция записи распространилась на указанное количество экземпляров
// С j: true MongoDB возвращает данные только после того,
// как запрошенное количество участников, включая основного, записали данные