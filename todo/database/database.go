package database

type Database struct {
	Conn string
	Test func() string
}

func PostgreConn() *Database {
	return &Database{Conn: "conexao.postgres.5432", Test: func() string { return "Conectado" }}
}
