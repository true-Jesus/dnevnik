package server

import (
	"database/sql"
	"dnevnik/repo"
	"dnevnik/usecases"
	"fmt"
	"net/http"
)

func RunServer(db *sql.DB) error {
	fmt.Println("zbc")
	//create repos
	repoSql := repo.NewRepo(db)
	//create usecases

	aufUc := usecases.NewAufUseCase(repoSql)
	dickUc := usecases.NewGradeUseCase(repoSql)

	useCases := NewUseCases(aufUc, dickUc)
	//create handlers
	h := NewHandlers(useCases)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", 8080), Handler: NewRouter(h)}
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server.Run #2: %s \n", err)
	}

	return nil
}
