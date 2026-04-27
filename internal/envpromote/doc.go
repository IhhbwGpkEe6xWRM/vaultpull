// Package envpromote implements promotion of secret maps between environments.
//
// A Promoter copies keys from a source secret map (e.g. staging) into a
// destination map (e.g. production). Callers may restrict which keys are
// eligible via an allow-list, and can enable dry-run mode to preview the
// operation without mutating the destination.
//
// Example:
//
//	p := envpromote.New(
//		envpromote.WithAllowList([]string{"DB_URL", "API_KEY"}),
//	)
//	res, err := p.Promote(stagingSecrets, productionSecrets)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("promoted %d keys\n", len(res.Promoted))
package envpromote
