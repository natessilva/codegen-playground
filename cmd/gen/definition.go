package main

type Definition struct {
	Services []Service
	Structs  []Struct
	Enums    []Enum
}

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name   string
	Input  string
	Output string
}

type Struct struct {
	Name   string
	Fields []Field
}

type Enum struct {
	Name   string
	Consts []Const
}

type Const struct {
	Name  string
	Value string
}

type Field struct {
	Name      string
	IsSlice   bool
	IsNumeric bool
	Type      string
	Tag       string
}
