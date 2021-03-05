package seeder

import (
	"database/sql"
	"fmt"

	"github.com/riskiramdan/evos/config"
)

// SeedUp seeding the database
func SeedUp() error {
	cfg, err := config.GetConfiguration()
	if err != nil {
		return fmt.Errorf("error when getting configuration: %s", err)
	}

	db, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		return fmt.Errorf("error when open postgres connection: %s", err)
	}
	defer db.Close()

	//Roles
	_, err = db.Exec(`
	INSERT INTO public."roles"
	(id, "name")
	VALUES(1, 'Admin');
	INSERT INTO public."roles"
	(id, "name")
	VALUES(2, 'Operator');
	INSERT INTO public."roles"
	(id, "name")
	VALUES(3, 'Guest');
	`)
	if err != nil {
		return err
	}

	//User
	_, err = db.Exec(`
	INSERT INTO public.users
	(id, "roleId", name, phone, "password", "token", "tokenExpiredAt")
	VALUES(9999, 1, 'admin', '082101010101', 'jLov', NULL, NULL);
	`)
	if err != nil {
		return err
	}

	//Character Type
	_, err = db.Exec(`
	INSERT INTO public."charactersType"
	(id, name, code)
	VALUES(1, 'Wizard', 1);
	INSERT INTO public."charactersType"
	(id, name, code)
	VALUES(2, 'Elf', 2);
	INSERT INTO public."charactersType"
	(id, name, code)
	VALUES(3, 'Hobbit', 3);
	`)
	if err != nil {
		return err
	}

	//Character
	_, err = db.Exec(`
	INSERT INTO public.characters
	(id, "characterTypeID", name, power)
	VALUES(9999, 1, 'Gandalf', 100);
	INSERT INTO public.characters
	(id, "characterTypeID", name, power)
	VALUES(9998, 2, 'Legolas', 60);
	INSERT INTO public.characters
	(id, "characterTypeID", name, power)
	VALUES(9997, 3, 'Frodo', 10);	
	`)
	if err != nil {
		return err
	}

	return nil
}
