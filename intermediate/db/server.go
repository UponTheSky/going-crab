package main

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Role struct {
	Id        int
	Name      string
	ActorName string
}

func main() {
	dbUrl := "postgres://gcrab:gcrab@localhost:5432/hangover"
	config, err := pgx.ParseConfig(dbUrl)

	if err != nil {
		log.Fatal(err)
	}

	connCtx := context.Background()
	conn, err := pgx.ConnectConfig(connCtx, config)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		closeCtx := context.Background()

		if err := conn.Close(closeCtx); err != nil {
			log.Fatal(err)
		}
	}()

	rows := []Role{
		{1, "Leslie Chow", "Ken Jeong"},
		{2, "Philip Wenneck", "Bradley Cooper"},
		{3, "Stuart Price", "Ed Helms"},
		{4, "Zach Galifianakis", "Alan Garner"},
		{5, "Teddy Srisai", "Mason Lee"},
	}
	txCtx := context.Background()

	tx, err := conn.Begin(txCtx)

	if err != nil {
		log.Fatalf("opening transaction error: %v", err)
	}

	defer func() {
		rbCtx := context.Background()
		tx.Rollback(rbCtx)
	}()

	execCtx := context.Background()

	for _, role := range rows {
		tag, err := conn.Exec(
			execCtx,
			"INSERT INTO general.role(id, name, actor_name) VALUES ($1, $2, $3)",
			role.Id, role.Name, role.ActorName,
		)

		if err != nil {
			log.Fatalf("executing insert operation error: %v", err)
		}

		if tag.RowsAffected() == 0 {
			log.Fatalf("inserting the current row has not been successful: %v", role)
		}
	}

	commitCtx := context.Background()

	if err = tx.Commit(commitCtx); err != nil {
		log.Fatalf("commit error: %v", err)
	}

	execCtx = context.Background()
	firstRow, err := conn.Query(execCtx, "SELECT id, name, actor_name FROM general.role WHERE id = $1", rows[0].Id)

	if err != nil {
		log.Fatal("sending the query and initializing the rows has not been successful", err)
	}

	row, err := pgx.CollectExactlyOneRow(firstRow, func(rawRow pgx.CollectableRow) (Role, error) {
		value, err := rawRow.Values()

		if err != nil || len(value) != 3 {
			return Role{}, err
		}

		id := value[0].(int32)
		name := value[1].(string)
		actorName := value[2].(string)

		return Role{int(id), name, actorName}, nil
	})

	if err != nil {
		log.Fatalf("either 0, or more than one rows exist: %v", err)
	}

	log.Println(row)

	// start transaction
	txCtx = context.Background()
	tx, err = conn.Begin(txCtx)

	if err != nil {
		log.Fatalf("opening transaction error: %v", err)
	}

	defer func() {
		rbCtx := context.Background()
		tx.Rollback(rbCtx)
	}()

	// update
	execCtx = context.Background()
	tag, err := conn.Exec(execCtx, "UPDATE general.role SET actor_name = $1 WHERE id = $2", strings.ToUpper(rows[0].ActorName), rows[0].Id)

	if err != nil {
		log.Fatalf("executing update operation error: %v", err)
	}

	if tag.RowsAffected() == 0 {
		log.Fatalf("row with id %v has not been updated", rows[0].Id)
	}

	// delete
	execCtx = context.Background()
	tag, err = conn.Exec(execCtx, "DELETE FROM general.role WHERE LOWER(name) LIKE $1", "%"+"teddy"+"%")

	if err != nil {
		log.Fatalf("executing delete operation error: %v", err)
	}

	if tag.RowsAffected() == 0 {
		log.Fatalln("row with name 'Teddy' has not been deleted")
	}

	// commit
	commitCtx = context.Background()

	if err = tx.Commit(commitCtx); err != nil {
		log.Fatalf("commit error: %v", err)
	}
}
