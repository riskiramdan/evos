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
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(1, 'Admin', '2021-02-04 17:22:00.028991', 'admin', '2021-02-04 17:22:00.028991', 'admin', NULL, NULL);
	INSERT INTO public."roles"
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(2, 'Operator', '2021-02-04 17:22:21.016457', 'admin', '2021-02-04 17:22:21.016457', 'admin', NULL, NULL);
	INSERT INTO public."roles"
	(id, "name", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(3, 'Guest', '2021-02-04 17:22:37.551864', 'admin', '2021-02-04 17:22:37.551864', 'admin', NULL, NULL);	
	`)
	if err != nil {
		return err
	}

	//User
	_, err = db.Exec(`
	INSERT INTO public.users
	(id, "roleId", name, phone, "password", "token", "tokenExpiredAt", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(9999, 1, 'admin', '082101010101', 'jLov', NULL, NULL, '2021-02-07 14:37:52.252246', 'admin', '2021-02-07 14:37:52.252246', 'admin', NULL, NULL);
	INSERT INTO public.users
	(id, "roleId", name, phone, "password", "token", "tokenExpiredAt", "createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy")
	VALUES(9998, 2, 'efishery', '082102020202', 'VWqV', NULL, NULL, '2021-02-07 14:38:07.292022', 'admin', '2021-02-07 14:38:07.292022', 'admin', NULL, NULL);	
	`)
	if err != nil {
		return err
	}

	return nil
}
