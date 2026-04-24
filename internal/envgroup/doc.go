// Package envgroup partitions a flat map of environment variables into named
// groups based on key-prefix rules.
//
// This is useful when a single Vault secret path contains keys for multiple
// logical services (e.g. DB_HOST, CACHE_URL) and each service needs its own
// scoped view of the configuration.
//
// Usage:
//
//	g, err := envgroup.New([]string{"db=DB", "cache=CACHE"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	groups := g.Split(secrets)
//	for _, grp := range groups {
//	    fmt.Println(grp.Name, grp.Values)
//	}
//
// Keys that do not match any rule are placed in the group with an empty name.
package envgroup
