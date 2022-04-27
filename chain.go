package dino

import (
	"reflect"
	"strings"
)

// DepLink describes a stack frame of currently called bindings.
type DepLink struct {
	ty      reflect.Type // Type requested from the container.
	binding Binding      // Binding used to realize the request.
}

// CyclicDependencyError occurs when a container cannot construct a service,
// because a transient service depends on itself.
//
// If cyclic dependencies were not checked, the application would fail with a stack overflow.
type CyclicDependencyError struct {
	chain []DepLink
}

func (e CyclicDependencyError) Error() string {
	return "Cannot satisfy cyclic dependency: " + formatChain(e.chain, true)
}

// formatChain describes a chain of dependencies as a human-readable string.
//
// If highlightLast is set, it will print all links matching the last one as uppercase.
func formatChain(chain []DepLink, highlightLast bool) string {

	// Store the last element of the chain, if one exists
	lastLink := DepLink{}
	if highlightLast && len(chain) > 0 {
		lastLink = chain[len(chain)-1]
	}

	var b strings.Builder
	for _, link := range chain {
		svcName := "???"
		if link.ty != nil {
			svcName = link.ty.String()
		}
		if highlightLast && link == lastLink {
			svcName = strings.ToUpper(svcName)
		}
		b.WriteString(svcName)
		b.WriteString(" (")
		typeName := "???"
		{
			bty := reflect.TypeOf(link.binding)
			if bty != nil {
				name := bty.Name()
				if name == "" && bty.Kind() == reflect.Pointer {
					name = bty.Elem().Name()
				}
				typeName = name
			}
		}
		b.WriteString(strings.TrimSuffix(typeName, "Binding"))
		b.WriteString(") ---> ")
	}

	return strings.TrimSuffix(b.String(), " ---> ")
}
