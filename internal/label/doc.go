// Package label implements key=value label tagging and selector-based filtering
// for secrets managed by vaultpull.
//
// Labels are arbitrary string pairs attached to secret metadata. A Tagger
// stamps a fixed set of labels onto metadata maps, while a Selector filters
// secret collections by requiring that each entry's metadata satisfies a given
// label Set.
//
// Typical usage:
//
//	tagger := label.New([]string{"env=prod", "team=platform"})
//	meta  := tagger.Apply(existingMeta)
//
//	filter := label.ParseFilter([]string{"env=prod"})
//	sel    := label.NewSelector(filter)
//	matched := sel.Filter(secrets, metaLookup)
package label
