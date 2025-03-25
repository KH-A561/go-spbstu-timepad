package postgres

type Config interface{}

type config struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func NewDefault() Config {
	return config{host: "localhost", port: 5432, user: "postgres", password: "pass", dbname: "mydb"}
}

func New(host string, port int, user string, password string, dbname string) Config {
	return config{host: host, port: port, user: user, password: password, dbname: dbname}
}
