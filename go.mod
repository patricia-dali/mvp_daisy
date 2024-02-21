module exemplo.com/server

go 1.21.4

replace exemplo.com/router => ./src

replace exemplo.com/login => ./src/login

replace exemplo.com/cadastro => ./src/cadastro

replace exemplo.com/index => ./src/index

replace exemplo.com/database => ./src/database

replace exemplo.com/logout => ./src/logout

replace exemplo.com/resetPassword => ./src/resetPassword

replace exemplo.com/openAI => ./src/openAI

replace exemplo.com/users => ./src/users

require (
	exemplo.com/database v0.0.0-00010101000000-000000000000
	exemplo.com/router v0.0.0-00010101000000-000000000000
	github.com/gorilla/sessions v1.2.2
	github.com/joho/godotenv v1.5.1
)

require (
	exemplo.com/cadastro v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/index v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/login v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/logout v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/openAI v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/resetPassword v0.0.0-00010101000000-000000000000 // indirect
	exemplo.com/users v0.0.0-00010101000000-000000000000 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/sashabaranov/go-openai v1.19.3 // indirect
	golang.org/x/crypto v0.18.0 // indirect
)
