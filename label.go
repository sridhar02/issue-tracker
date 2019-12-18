package main

import (
	"database/sql"
	"fmt"
)

type Label struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

func GetLabel(db *sql.DB, Id int) (Label, error) {

	var name, description, color string

	row := db.QueryRow("SELECT  name,description,color FROM labels WHERE id=$1", Id)
	err := row.Scan(&name, &description, &color)
	if err != nil {
		return Label{}, err
	}

	label := Label{
		ID:          Id,
		Name:        name,
		Description: description,
		Color:       color,
	}

	return label, nil
}

func CreateLabel(db *sql.DB, label Label) error {

	_, err := db.Exec(`INSERT INTO labels ( name, description,color) VALUES ( $1, $2, $3 ) `, label.Name, label.Description, label.Color)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func DeleteLabel(db *sql.DB, id int) error {

	_, err := db.Exec(`DELETE FROM labels WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}


