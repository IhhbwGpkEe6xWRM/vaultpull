// Package envalias provides key aliasing for vault secrets.
//
// It allows operators to rename vault secret keys to environment variable
// names that differ from the vault path keys. Aliases are specified as
// "from=to" pairs. When applied, the From key is removed and its value
// is written under the To key, unless To already exists.
//
// Example usage:
//
//	m := envalias.New([]string{"db_password=DATABASE_PASSWORD"})
//	result := m.Apply(secrets)
package envalias
