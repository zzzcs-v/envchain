// Package inject merges one or more named variable sources into a destination
// map. It supports optional key prefixing and overwrite control, making it
// suitable for layering environment configs across contexts.
//
// Basic usage:
//
//	inj := inject.New(inject.Options{Overwrite: false, Prefix: "APP_"})
//	results, err := inj.Inject(dst, sources)
//	for _, r := range results {
//		fmt.Printf("%s: injected=%d skipped=%d\n", r.Source, r.Injected, r.Skipped)
//	}
package inject
