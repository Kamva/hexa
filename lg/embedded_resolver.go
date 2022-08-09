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

func NewEmbeddedResolver(packages ...*Package) *EmbeddedResolver {
	return NewEmbeddedResolverWithOpts(packages, nil)
}

func NewEmbeddedResolverWithOpts(packages []*Package, filter EmbeddedResolveFilter) *EmbeddedResolver {
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

func (r *EmbeddedResolver) resolveIfaceMethods(ifacePkg *Package, f *File, iface *Interface, em *EmbeddedField) error {
	if (r.filter != nil && r.filter(em)) || em.IsResolved {
		return nil
	}

	fieldPkg := ifacePkg // by default, we expect embedded field be in the current package.
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

	methods := embeddedIface.Methods
	if !IsSamePackage(fieldPkg, ifacePkg) {
		methods = UseMethodsInPackage(fieldPkg, methods)
	}

	iface.Methods = append(iface.Methods, methods...)
	em.IsResolved = true
	return nil
}

func (r *EmbeddedResolver) resolveStructFields(structPkg *Package, f *File, strct *Struct, em *EmbeddedField) error {
	if (r.filter != nil && r.filter(em)) || em.IsResolved {
		return nil
	}

	fieldPkg := structPkg // by default, we expect embedded field be in the current package.
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

	fields := embeddedStruct.Fields
	if !IsSamePackage(fieldPkg, structPkg) {
		fields = UseFieldsInPackage(fieldPkg, fields)
	}

	strct.Fields = append(strct.Fields, fields...)
	em.IsResolved = true
	return nil
}

// UseInterfaceInPackage updates the interface to be able to use it in another package.
func UseInterfaceInPackage(from *Package, iface *Interface) *Interface {
	return &Interface{
		Doc:         iface.Doc,
		Annotations: iface.Annotations,
		Name:        iface.Name,
		Embedded:    UseEmbeddedFieldsInPackage(from, iface.Embedded),
		Methods:     UseMethodsInPackage(from, iface.Methods),
	}
}

// UseStructInPackage updates the struct to be able to use it in another package.
func UseStructInPackage(from *Package, strct *Struct) *Struct {
	return &Struct{
		Doc:         strct.Doc,
		Annotations: strct.Annotations,
		Name:        strct.Name,
		Embedded:    UseEmbeddedFieldsInPackage(from, strct.Embedded),
		Fields:      UseFieldsInPackage(from, strct.Fields),
	}
}

// UseMethodsInPackage updates the method's params and results to use in another package.
// e.g., when want to use checkHealth(h Health) to another package, it should be checkHealth(h hexa.Health).
func UseMethodsInPackage(from *Package, methods []*Method) []*Method {
	l := make([]*Method, len(methods))
	for i, m := range methods {
		params := make([]*MethodParam, len(m.Params))
		results := make([]*MethodResult, len(m.Results))

		// add the package's name of the "from" package to methods params and results in it.
		// e.g, converts `hi(h Health)` to `hi(h hexa.Health)` to use in non-hexa packages.
		for i, p := range m.Params {
			params[i] = &MethodParam{
				Name: p.Name,
				Type: UseTypeInPackage(from, p.Type),
			}
		}

		for i, r := range m.Results {
			results[i] = &MethodResult{
				Name: r.Name,
				Type: UseTypeInPackage(from, r.Type),
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

// UseEmbeddedFieldsInPackage updates the embedded fields to be able to use in another package.
func UseEmbeddedFieldsInPackage(from *Package, fields []*EmbeddedField) []*EmbeddedField {
	l := make([]*EmbeddedField, len(fields))

	for i, field := range fields {
		l[i] = &EmbeddedField{
			IsResolved:  field.IsResolved,
			Doc:         field.Doc,
			Annotations: field.Annotations,
			Type:        UseTypeInPackage(from, field.Type),
			Tag:         field.Tag,
		}
	}

	return fields
}

// UseFieldsInPackage updates fields to be able to use in the another package.
// e.g., when we want to use
// ```
// type Hi struct{h Health}
// ```
// in another package, it should be:
// ```
//type Hi struct{h hexa.Health}
// ```
func UseFieldsInPackage(from *Package, fields []*Field) []*Field {
	l := make([]*Field, len(fields))

	// add the package's name of the "from" package to fields.
	// e.g, converts field `h Health` to `h hexa.Health` to use in non-hexa packages.
	for i, field := range fields {
		l[i] = &Field{
			Doc:         field.Doc,
			Annotations: field.Annotations,
			Name:        field.Name,
			Type:        UseTypeInPackage(from, field.Type),
			Tag:         field.Tag,
		}
	}

	return l
}

// UseTypeInPackage returns the type that we can use in another package.
// e.g., if we provide Health from the hexa package. it returns hexa.Health.
func UseTypeInPackage(from *Package, t string) string {
	parsedPkg, _ := parseType(t)

	// actually we use isPrivateType to check if it's a primitive type, because all
	// primitive types start with a lower case just like a private type.
	if parsedPkg == "" && !isPrivateType(t) {
		return SetPackageOnType(from.Name, t)
	}
	return t
}
