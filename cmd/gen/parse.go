package main

import (
	"fmt"
	"go/types"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func ParseDefinition(pattern string) (Definition, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}
	result := Definition{}
	pattern, err := filepath.Abs(pattern)
	if err != nil {
		return result, err
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return result, errors.Wrap(err, "load")
	}
	if len(pkgs) != 1 {
		return result, fmt.Errorf("pattern %s must contain 1 package", pattern)
	}
	pkg := pkgs[0]
	if len(pkg.Imports) != 0 {
		return result, errors.New("package must not import")
	}
	scope := pkg.Types.Scope()
	seenNames := make(map[string]bool)
	for _, name := range scope.Names() {
		seenNames[name] = false
		obj := scope.Lookup(name)

		if !obj.Exported() {
			return result, fmt.Errorf("object %s must be exported", name)
		}

		if i, ok := obj.Type().Underlying().(*types.Interface); ok {
			service, err := parseService(i, name)
			if err != nil {
				return result, errors.Wrapf(err, "service %s failed to parse", name)
			}
			seenNames[name] = true
			result.Services = append(result.Services, service)
		}
	}
	structMap := make(map[string]Struct)
	for _, s := range result.Services {
		for _, m := range s.Methods {
			err := parseStruct(structMap, scope.Lookup(m.Input).Type().Underlying().(*types.Struct), m.Input)
			if err != nil {
				return result, errors.Wrapf(err, "struct %s failed to parse", m.Input)
			}
			err = parseStruct(structMap, scope.Lookup(m.Output).Type().Underlying().(*types.Struct), m.Output)
			if err != nil {
				return result, errors.Wrapf(err, "struct %s failed to parse", m.Output)
			}
		}
	}
	for _, v := range structMap {
		seenNames[v.Name] = true
		result.Structs = append(result.Structs, v)
	}
	for name, seen := range seenNames {
		if !seen {
			return result, fmt.Errorf("scope %s must be referenced by a service", name)
		}
	}
	return result, nil
}

func parseService(s *types.Interface, name string) (Service, error) {
	svc := Service{
		Name: name,
	}
	for i := 0; i < s.NumMethods(); i++ {
		m := s.Method(i)
		if !m.Exported() {
			return svc, fmt.Errorf("method %s must be exported", m.Name())
		}
		sig := m.Type().(*types.Signature)
		params := sig.Params()
		if params.Len() != 1 {
			return svc, fmt.Errorf("method %s must have 1 parameter", m.Name())
		}
		input := params.At(0)
		inputName, ok := input.Type().(*types.Named)
		if !ok {
			return svc, fmt.Errorf("method %s must have a named struct type as its input", m.Name())
		}
		if _, ok := input.Type().Underlying().(*types.Struct); !ok {
			return svc, fmt.Errorf("method %s must have a named struct type as its input", m.Name())
		}
		results := sig.Results()
		if results.Len() != 1 {
			return svc, fmt.Errorf("method %s must have 1 output", m.Name())
		}
		output := results.At(0)
		outputName, ok := output.Type().(*types.Named)
		if !ok {
			return svc, fmt.Errorf("method %s must have a named struct type as its output", m.Name())
		}
		if _, ok := output.Type().Underlying().(*types.Struct); !ok {
			return svc, fmt.Errorf("method %s must have a named struct type as its output", m.Name())
		}
		method := Method{
			Name:   m.Name(),
			Input:  inputName.Obj().Name(),
			Output: outputName.Obj().Name(),
		}
		svc.Methods = append(svc.Methods, method)
	}
	return svc, nil
}

func parseStruct(structMap map[string]Struct, s *types.Struct, name string) error {
	if _, ok := structMap[name]; ok {
		return nil
	}
	st := Struct{
		Name: name,
	}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		if !f.Exported() {
			return fmt.Errorf("struct field %s must be exported", f.Name())
		}
		field := Field{
			Name: f.Name(),
			Tag:  s.Tag(i),
		}
		t := f.Type()
		if i, ok := t.(*types.Slice); ok {
			field.IsSlice = true
			t = i.Elem()
		}
		switch ut := t.Underlying().(type) {
		case *types.Basic:
			if !isConstType(ut) {
				return fmt.Errorf("struct field %s must be a const basic type", field.Name)
			}
			field.IsNumeric = IsNumeric(ut)
			break
		case *types.Struct:
			n, ok := t.(*types.Named)
			if !ok {
				return fmt.Errorf("struct field %s must be a named struct type", field.Name)
			}
			err := parseStruct(structMap, ut, n.Obj().Name())
			if err != nil {
				return errors.Wrapf(err, "struct %s for field %s failed to parse", n.Obj().Name(), field.Name)
			}
			break
		default:
			return fmt.Errorf("struct field %s must be a const basic or a name struct type", field.Name)
		}
		field.Type = types.TypeString(t, func(*types.Package) string { return "" })
		st.Fields = append(st.Fields, field)
	}
	structMap[name] = st
	return nil
}

func isConstType(t *types.Basic) bool {
	return t.Info()&types.IsConstType != 0
}

func IsNumeric(t *types.Basic) bool {
	return t.Info()&types.IsNumeric != 0
}
