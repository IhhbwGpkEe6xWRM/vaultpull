// Package envflow implements a sequential secret-transformation pipeline.
//
// A Pipeline is built from an ordered slice of named Stage values. Each Stage
// receives the output of the previous stage as its input, allowing complex
// processing graphs to be composed from small, testable functions.
//
// Example usage:
//
//	p, err := envflow.New([]envflow.Stage{
//		{Name: "uppercase", Apply: transform.UppercaseKeys},
//		{Name: "prefix",    Apply: myPrefixFunc},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, results, err := p.Run(secrets)
//
// Run stops at the first stage that returns an error and reports partial
// results so callers can log which stages succeeded before the failure.
package envflow
