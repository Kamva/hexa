package lg

import (
	"fmt"

	"github.com/kamva/tracer"
)

type EmbeddedResolveFilter func(e *EmbeddedField) bool
type EmbeddedResolver struct {
	packages []*Package
	pm       map[string]*Package // map's key is the package's path. e.g., github.com/kamva/hexa
	filter   EmbeddedResolveFilter
}

func NewEmbeddedResolver(filter EmbeddedResolveFilter, packages ...*Package) *EmbeddedResolver {
	pm := make(map[string]*Package)
	for _, p := range packages {
		pm[p.Path] = p
	}

	return &EmbeddedResolver{
		packages: packages,
		pm:       pm,
		filter:   filter,
	}
}

func (r *EmbeddedResolver) Resolve() error {
	for _, p := range r.packages {
		for _, f := range p.Files {
			// resolve embedded fields in all interfaces
			for _, iface := range f.Interfaces {
				for _, em := range iface.Embedded {
					if err := r.resolveIfaceMethods(p, f, iface, em); err != nil {
						return tracer.Trace(err)
					}
				}
			}

			// resolve embedded fields in all structs
			for _, strct := range f.Structs {
				for _, em := range strct.Embedded {
					if err := r.resolveStructFields(p, f, strct, em); err != nil {
						return tracer.Trace(err)
					}
				}
			}
		}
	}
	return nil
}

func (r *EmbeddedResolver) resolveIfaceMethods(p *Package, f *File, iface *Interface, em *EmbeddedField) error {
	if (r.filter != nil && r.filter(em)) || em.IsResolved {
		return nil
	}

	fieldPkg := p // by default, we expect embedded field be in the current package.
	pname, tname := parseType(em.Type)
	if pname != "" { // Find package of the embedded field.
		if fieldPkg = r.pm[f.ImportMap[pname]]; fieldPkg == nil {
			return tracer.Trace(fmt.Errorf("package with path %s is not parsed, add it to your parse litst please", f.ImportMap[pname]))
		}
	}

	embeddedIfaceFile, embeddedIface := fieldPkg.FindInterface(tname)
	if embeddedIface == nil {
		return tracer.Trace(fmt.Errorf("can not resolve embedded interface, interface with name %s in the package: %s not found", tname, pname))
	}

	for _, em := range embeddedIface.Embedded {
		if err := r.resolveIfaceMethods(fieldPkg, embeddedIfaceFile, embeddedIface, em); err != nil {
			return tracer.Trace(err)
		}
	}

	iface.Methods = append(iface.Methods, prepareMethodsToUseInPackage(fieldPkg, p, embeddedIface.Methods)...)
	em.IsResolved = true
	return nil
}

func (r *EmbeddedResolver) resolveStructFields(p *Package, f *File, strct *Struct, em *EmbeddedField) error {
	if r.filter(em) || em.IsResolved {
		return nil
	}

	fieldPkg := p // by default, we expect embedded field be in the current package.
	pname, tname := parseType(em.Type)
	if pname != "" {
		if fieldPkg = r.pm[f.ImportMap[pname]]; fieldPkg == nil {
			return tracer.Trace(fmt.Errorf("package with path %s is not parsed, add it to your parse litst please", f.ImportMap[pname]))
		}
	}

	// If the embedded field is an interface, we can skip it.
	if _, iface := fieldPkg.FindInterface(tname); iface != nil {
		return nil
	}

	embeddedStructFile, embeddedStruct := fieldPkg.FindStruct(tname)
	if embeddedStruct == nil {
		return tracer.Trace(fmt.Errorf("can not resolve embedded struct, struct with name %s in the package: %s not found", tname, pname))
	}

	for _, em := range embeddedStruct.Embedded {
		if err := r.resolveStructFields(fieldPkg, embeddedStructFile, embeddedStruct, em); err != nil {
			return tracer.Trace(err)
		}
	}

	strct.Fields = append(strct.Fields, prepareFieldsToUseInPackage(fieldPkg, p, embeddedStruct.Fields)...)
	em.IsResolved = true
	return nil
}

// prepareMethodsToUseInPackage updates the method's params and results to use in the the "to" package.
// e.g., when want to use checkHealth(h Health) to another package, it should be checkHealth(h hexa.Health).
func prepareMethodsToUseInPackage(from, to *Package, methods []*Method) []*Method {
	if from == to {
		return methods
	}

	l := make([]*Method, len(methods))
	for i, m := range methods {
		params := make([]*MethodParam, len(m.Params))
		results := make([]*MethodResult, len(m.Results))

		// add the package's name of the "from" package to methods params and results in it.
		// e.g, converts `hi(h Health)` to `hi(h hexa.Health)` to use in non-hexa packages.
		for i, p := range m.Params {
			paramPkg, paramType := parseType(p.Type)
			if paramPkg == "" {
				paramPkg = from.Name
			}
			params[i] = &MethodParam{
				Name: p.Name,
				Type: fmt.Sprintf("%s.%s", paramPkg, paramType),
			}
		}

		for i, r := range m.Results {
			resultPkg, resultType := parseType(r.Type)
			if resultPkg == "" {
				resultPkg = from.Name
			}

			results[i] = &MethodResult{
				Name: r.Name,
				Type: constructType(resultPkg, resultType),
			}
		}

		l[i] = &Method{
			Doc:         m.Doc,
			Annotations: m.Annotations,
			Name:        m.Name,
			Params:      params,
			Results:     results,
		}
	}

	return l
}

// prepareFieldsToUseInPackage updates fields to be able to use in the "to" package.
// e.g., when we want to use
// ```
// type Hi struct{h Health}
// ```
// in another package, it should be:
// ```
//type Hi struct{h hexa.Health}
// ```
func prepareFieldsToUseInPackage(from, to *Package, fields []*Field) []*Field {
	if from == to {
		return fields
	}

	l := make([]*Field, len(fields))

	// add the package's name of the "from" package to fields.
	// e.g, converts field `h Health` to `h hexa.Health` to use in non-hexa packages.
	for i, field := range fields {
		fieldPkg, fieldType := parseType(field.Type)
		if fieldPkg == "" {
			fieldPkg = from.Name
		}

		l[i] = &Field{
			Doc:         field.Doc,
			Annotations: field.Annotations,
			Name:        field.Name,
			Type:        constructType(fieldPkg, fieldType),
			Tag:         field.Tag,
		}
	}

	return l
}
